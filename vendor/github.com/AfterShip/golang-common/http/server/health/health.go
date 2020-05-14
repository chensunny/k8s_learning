package health

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

const requestPath = "/devops/status"

// NewStatus creates a Status with options.
func NewStatus(opts ...Option) *Status {
	status := &Status{
		status: http.StatusServiceUnavailable,
	}
	for _, opt := range opts {
		opt(status)
	}
	return status
}

type Option func(status *Status)

func DefaultStatus(statusCode int32) Option {
	return func(status *Status) {
		status.status = statusCode
	}
}

func AllowRemoteUpdate(allowed bool) Option {
	return func(status *Status) {
		status.allowRemoteUpdate = allowed
	}
}

type Status struct {
	status            int32
	allowRemoteUpdate bool
}

func (s *Status) Online() {
	atomic.StoreInt32(&s.status, int32(http.StatusOK))
}

func (s *Status) Offline() {
	atomic.StoreInt32(&s.status, int32(http.StatusServiceUnavailable))
}

func (s *Status) RequestURI() string {
	return requestPath
}

func (s Status) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(s.StatusCode())
	fmt.Fprintln(w, s.StatusText())
}

func (s Status) StatusCode() int {
	return int(s.status)
}

func (s Status) StatusText() string {
	switch s.StatusCode() {
	case http.StatusOK:
		return http.StatusText(http.StatusOK)
	case http.StatusServiceUnavailable:
		fallthrough
	default:
		return http.StatusText(http.StatusServiceUnavailable)
	}
}

func (s *Status) RegisterHttpHandler(mux *http.ServeMux) {
	mux.HandleFunc(requestPath, s.HttpHandlerFunc())
}

func (s Status) isLocalIP(ip string) bool {
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		addr, _, _ := net.ParseCIDR(a.String())
		if net.ParseIP(ip).Equal(addr) {
			return true
		}
	}
	return false
}

// HttpHandler creates a new HTTP handler.
func (s *Status) HttpHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodHead, http.MethodGet:
			s.WriteResponse(w)
		case http.MethodPut, http.MethodPost:
			remoteIP := strings.Split(r.RemoteAddr, ":")[0]
			if !s.allowRemoteUpdate && !s.isLocalIP(remoteIP) {
				http.Error(w, "status remote update is not allowed", http.StatusBadRequest)
				return
			}

			if sParam := r.FormValue("status"); sParam != "" {
				if sInt, err := strconv.ParseInt(sParam, 10, 64); err == nil {
					switch sInt {
					case http.StatusOK:
						s.Online()
					case http.StatusServiceUnavailable:
						s.Offline()
					default:
						http.Error(w, "Invalid Service Status: "+sParam, http.StatusBadRequest)
						return
					}

					w.WriteHeader(http.StatusOK)
					fmt.Fprintf(w, "Service Status Updated, Current Status: %d", sInt)
					return
				}
			}

			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}
