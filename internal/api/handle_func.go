package api

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mtavano/admoai-takehome/internal/metrics"
)

type handler func(*gin.Context, *Context) (any, int, error)

func HandleFunc(fn handler, ctx *Context) func(*gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.NewString()
		
		log.Printf("request received %s", requestID)
		
		payload, statusCode, err := fn(c, ctx)
		
		elapsed := time.Since(start)
		
		// Record metrics using Prometheus collector
		collector := metrics.GetCollector()
		if collector != nil {
			endpoint := c.FullPath()
			if endpoint == "" {
				endpoint = c.Request.URL.Path
			}
			status := strconv.Itoa(statusCode)
			collector.RecordHTTPRequest(c.Request.Method, endpoint, status, elapsed)
		}

		if err != nil {
			log.Printf("request error %s %v", requestID, elapsed)
			c.JSON(statusCode, map[string]any{
				"error": err.Error(),
			})
			return
		}

		log.Printf("request responded %s %v", requestID, elapsed)
		
		// Don't send response if it's already been sent (like in metrics handler)
		if !c.Writer.Written() {
			c.JSON(statusCode, payload)
		}
	}
}
