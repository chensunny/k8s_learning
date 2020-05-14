package tracing

import (
	"context"

	"github.com/AfterShip/golang-common/errors"
	newrelic "github.com/newrelic/go-agent"
)

func GetTraceIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(ContextKeyTraceID).(string)
	if ok {
		return requestID
	} else {
		return ""
	}
}

func GetCloudflareRayFromContext(ctx context.Context) string {
	ray, ok := ctx.Value(ContextKeyCloudflareRay).(string)
	if ok {
		return ray
	} else {
		return ""
	}
}

func GetRequestMethodFromContext(ctx context.Context) string {
	requestMethod, ok := ctx.Value(ContextKeyRequestMethod).(string)
	if ok {
		return requestMethod
	} else {
		return ""
	}
}

func GetRequestPathFromContext(ctx context.Context) string {
	requestPath, ok := ctx.Value(ContextKeyRequestPath).(string)
	if ok {
		return requestPath
	} else {
		return ""
	}
}

func GetTaskProcessErrorFromContext(ctx context.Context) *errors.APIError {
	err, ok := ctx.Value(ContextKeyTaskProcessError).(*errors.APIError)
	if ok {
		return err
	} else {
		return nil
	}
}

func ContextWithTraceID(parent context.Context, traceID string) context.Context {
	return context.WithValue(parent, ContextKeyTraceID, traceID)
}

func ContextWithCloudflareRay(parent context.Context, cloudflareRay string) context.Context {
	return context.WithValue(parent, ContextKeyCloudflareRay, cloudflareRay)
}

func ContextWithNewrelicTxn(parent context.Context, txn newrelic.Transaction) context.Context {
	return context.WithValue(parent, NewrelicGinTransactionContextKey, txn)
}

func ContextWithRequestMethod(parent context.Context, requestMethod string) context.Context {
	return context.WithValue(parent, ContextKeyRequestMethod, requestMethod)
}

func ContextWithRequestPath(parent context.Context, requestPath string) context.Context {
	return context.WithValue(parent, ContextKeyRequestPath, requestPath)
}

func ContextWithTaskProcessError(parent context.Context, err *errors.APIError) context.Context {
	return context.WithValue(parent, ContextKeyTaskProcessError, err)
}
