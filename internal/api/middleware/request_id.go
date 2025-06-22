package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestID struct{}

func NewRequestID() *RequestID {
	return &RequestID{}
}

func (rid *RequestID) Setup(engine *gin.Engine) {
	engine.Use(rid.handler())
}

func (rid *RequestID) handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		// Set request ID in gin context
		c.Set("request_id", requestID)

		// Set request ID in response header
		c.Writer.Header().Set("X-Request-ID", requestID)

		// Add request ID to the request context
		ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
