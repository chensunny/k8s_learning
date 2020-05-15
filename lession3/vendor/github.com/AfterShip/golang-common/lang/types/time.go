package types

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"time"
)

const (
	rfc3339         = "2006-01-02T15:04:05-0700"
	rfc3339Colon    = "2006-01-02T15:04:05-07:00"
	rfc3339WithZulu = "2006-01-02T15:04:05Z0700"
	rfc3339Slash    = "2006/01/02T15:04:05-0700"
	custom1         = "2006/01/02 15:04:05"
)

var (
	_              WrappedType = Time{}
	TimeType                   = reflect.TypeOf(Time{})
	nativeTimeType             = reflect.TypeOf(Time{}.val)

	allowedTimeFormats = []string{rfc3339, rfc3339Colon, rfc3339WithZulu, time.RFC822, time.RFC822Z,
		time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123,
		time.RFC1123Z, time.RFC3339, time.RFC3339Nano, time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro,
		time.StampNano, rfc3339Slash, custom1,
	}
)

// swagger:type time
type Time struct {
	val        time.Time
	isAssigned bool
	isNull     bool
}

// types.Time constructor.
// status[0]: isAssigned, whether field exists in JSON payloads or not.
// status[1]: isNull, whether field is JSON-Null / Spanner-Null / any other null value or not
func NewTime(val time.Time, status ...bool) Time {
	switch len(status) {
	case 0:
		return newTime(val, true, false)
	case 1:
		if !status[0] {
			return newTime(time.Time{}, false, false)
		}
		return newTime(val, status[0], false)
	default:
		if !status[0] {
			return newTime(time.Time{}, false, false)
		}
		if status[1] {
			return newTime(time.Time{}, true, true)
		}
		return newTime(val, status[0], status[1])
	}
}

func NewNullTime() Time {
	return newTime(time.Time{}, true, true)
}

func newTime(val time.Time, assigned bool, null bool) Time {
	return Time{
		val:        val,
		isAssigned: assigned,
		isNull:     null,
	}
}

func (n Time) Val() interface{} {
	return n.val
}

func (n Time) Time() time.Time {
	return n.Val().(time.Time)
}

func (n Time) Assigned() bool {
	return n.isAssigned
}

func (n Time) Null() bool {
	return n.isNull
}

func (n Time) Type() reflect.Type {
	return TimeType
}

func (n Time) NativeType() reflect.Type {
	return nativeTimeType
}

func (n *Time) SetToUTC() {
	n.val = n.val.UTC()
}

func (n *Time) CopyFrom(from Time, options ...Option) {
	n.val = from.Time()
	n.isAssigned = from.Assigned()
	n.isNull = from.Null()

	if len(options) > 0 {
		for _, opt := range options {
			opt(n)
		}
	}
}

func (n *Time) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, jsonNull) {
		n.isAssigned = true
		n.isNull = true
		return nil
	}
	if len(data) > 0 {
		timeStr, err := strconv.Unquote(string(data))
		if err != nil {
			return err
		}
		if len(timeStr) == 0 {
			n.isAssigned = false
			n.isNull = false
			return nil
		}
		var errs []error
		for _, format := range allowedTimeFormats {
			tmp, err := time.Parse(format, timeStr)
			if err == nil {
				n.val = tmp
				errs = []error{}
				break
			}
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			return errs[0]
		}
		n.isAssigned = true
		n.isNull = false
	}
	return nil
}

func (n Time) MarshalJSON() ([]byte, error) {
	if n.isAssigned {
		if n.isNull {
			return jsonNull, nil
		}
		return json.Marshal(n.val.Format(rfc3339Colon))
	}
	return json.Marshal("")
}
