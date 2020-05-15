package tracing

const (
	HeaderTraceID    = "am-trace-id"
	HeaderRequestID  = "Request-ID"
	HeaderXRequestID = "X-Request-ID"
	//gcp LB trace context header key , can't be changed
	HeaderXCloudTraceContext = "x-cloud-trace-context"
	HeaderCloudflareRay      = "CF-Ray"

	ContextKeyTraceID       = "automizelyTraceID"
	ContextKeyCloudflareRay = "cloudflareRay"

	ContextKeyRequestMethod = "requestMethod"
	ContextKeyRequestPath   = "requestPath"

	//from newrelic go-agent library
	NewrelicGinTransactionContextKey = "newRelicTransaction"

	ContextKeyTaskProcessError = "taskProcessError"
)
