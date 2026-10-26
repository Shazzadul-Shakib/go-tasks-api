package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

// RequestIDKey is the gin context key (and log field) holding the request ID.
const RequestIDKey = "request_id"

const requestIDHeader = "X-Request-ID"

// RequestID reuses the incoming X-Request-ID header or generates one, stores it
// in the gin context and echoes it back on the response.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" {
			id = newRequestID()
		}
		c.Set(RequestIDKey, id)
		c.Header(requestIDHeader, id)
		c.Next()
	}
}

func newRequestID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
