package api

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handler func(*gin.Context, *Context) (any, int, error)

func HandleFunc(fn handler, ctx *Context) func(*gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.NewString()
		
		log.Printf("request received %s", requestID)
		
		payload, statusCode, err := fn(c, ctx)
		
		elapsed := time.Since(start)

		if err != nil {
			log.Printf("request error %s %v", requestID, elapsed)
			c.JSON(statusCode, map[string]any{
				"error": err.Error(),
			})
			return
		}

		log.Printf("request responded %s %v", requestID, elapsed)
		c.JSON(statusCode, payload)
	}
}
