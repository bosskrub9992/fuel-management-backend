package middlewares

type (
	ContextKey  string
	LogFieldKey string
)

const (
	ContextKeyRequestID ContextKey = "request_id"
	ContextKeyHTMXData  ContextKey = "htmx_data"
)

const (
	LogFieldKeyReqID LogFieldKey = "reqId"
)
