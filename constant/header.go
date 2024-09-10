package constant

import "net/textproto"

var (
	XServiceName = textproto.CanonicalMIMEHeaderKey("x-service-name")
	XApiKey      = textproto.CanonicalMIMEHeaderKey("x-api-key")
	XRequestAt   = textproto.CanonicalMIMEHeaderKey("x-request-at")
	XRequestID   = textproto.CanonicalMIMEHeaderKey("x-request-id")
)
