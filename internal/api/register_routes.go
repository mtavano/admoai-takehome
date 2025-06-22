package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mtavano/admoai-takehome/internal/store"
)

type Context struct {
	Db store.Database
}

func RegisterRoutes(ctx *Context, engine *gin.Engine) {
	// TODO: add cors for frontend
	//corsMiddleware := middleware.NewCors()
	//requestIDMiddleware := middleware.NewRequestID()

	// Setup middlewares
	//corsMiddleware.Setup(engine, nil)
	//requestIDMiddleware.Setup(engine)

	engine.GET("/health", HandleFunc(func(c *gin.Context, ctx *Context) (any, int, error) {
		return map[string]any{
			"status": "running",
			//"request_id": c.GetString("request_id"),
		}, http.StatusOK, nil
	}, ctx))

	v1Router := engine.Group("/v1")

	v1Router.POST("/ads", HandleFunc(PostAdsHandler, ctx))
}
