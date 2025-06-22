package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mtavano/admoai-takehome/internal/api/middleware"
	"github.com/mtavano/admoai-takehome/internal/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Context struct {
	Db store.Database
}

func RegisterRoutes(ctx *Context, engine *gin.Engine) {
	corsMiddleware := middleware.NewCors()
	requestIDMiddleware := middleware.NewRequestID()

	// Setup middlewares
	corsMiddleware.Setup(engine, nil)
	requestIDMiddleware.Setup(engine)

	engine.GET("/health", HandleFunc(func(c *gin.Context, ctx *Context) (any, int, error) {
		return map[string]any{
			"status":     "running",
			"request_id": c.GetString("request_id"),
		}, http.StatusOK, nil
	}, ctx))

	// Metrics endpoint for Prometheus
	engine.GET("/metrics", HandleFunc(MetricsHandler, ctx))

	v1Router := engine.Group("/v1")

	v1Router.POST("/ads", HandleFunc(PostAdsHandler, ctx))
	v1Router.GET("/ads/:id", HandleFunc(GetAdsByIDHandler, ctx))
	v1Router.GET("/ads", HandleFunc(GetAdsByFiltersHandler, ctx))
	v1Router.POST("/ads/:id/deactivate", HandleFunc(PostDeactivateAdsHandler, ctx))
}

// MetricsHandler handles the /metrics endpoint
func MetricsHandler(c *gin.Context, ctx *Context) (any, int, error) {
	// Update current ad counts from database
	updateAdCounts(ctx)
	
	// Use Prometheus HTTP handler
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	
	// Return nil since promhttp.Handler handles the response
	return nil, http.StatusOK, nil
}

// updateAdCounts fetches current ad counts from database
func updateAdCounts(ctx *Context) {
	// This will be implemented when we add the metrics package
	// For now, we'll leave it empty since the metrics are handled by the Prometheus collector
}
