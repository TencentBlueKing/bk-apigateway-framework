package common

const (
	// RequestIDCtxKey Request ID 在 context 中的 key
	RequestIDCtxKey = "requestID"

	// RequestIDHeaderKey Request ID 在 HTTP Header 中的 key
	RequestIDHeaderKey = "X-Request-Id"

	// RequestIDLogKey Request ID 在日志中的 key
	RequestIDLogKey = RequestIDCtxKey

	// TraceIDLogKey trace id 在日志中的 key
	TraceIDLogKey = "otelTraceID"

	// SpanIDLogKey span id 在日志中的 key
	SpanIDLogKey = "otelSpanID"

	// ErrorCtxKey error 在 context 中的 key
	ErrorCtxKey = "error"

	// UserIDKey user id 在 cookies / session 中的 key
	UserIDKey = "bk_uid"

	// UserTokenKey user token 在 cookies / session 中的 key
	UserTokenKey = "bk_ticket"

	// UserLanguageKey user language 在 cookies / session 中的 key
	UserLanguageKey = "blueking_language"
)
