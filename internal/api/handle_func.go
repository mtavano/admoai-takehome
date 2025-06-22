package api

import (
	"log"

	"github.com/gin-gonic/gin"
)

type handler func(*gin.Context, *Context) (any, int, error)

func HandleFunc(fn handler, ctx *Context) func(*gin.Context) {
	return func(c *gin.Context) {
		//start := time.Now()
		payload, statusCode, err := fn(c, ctx)
		log.Println("api: call received")

		if err != nil {
			log.Println("api: Handler error", err)
			c.JSON(statusCode, map[string]any{
				"error": err.Error(),
			})
			return
		}

		c.JSON(statusCode, payload)
	}
}
