package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mtavano/admoai-takehome/internal/metrics"
	"github.com/mtavano/admoai-takehome/internal/store"
	"github.com/mtavano/admoai-takehome/internal/store/query"
	"github.com/pkg/errors"
)

func PostDeactivateAdsHandler(c *gin.Context, ctx *Context) (any, int, error) {
	// Get ID from path parameters
	id := c.Param("id")
	if id == "" {
		return gin.H{
			"error": "ID parameter is required",
		}, http.StatusBadRequest, nil
	}

	// Set status to inactive
	status := store.AdvertiseStatusInactive

	// Create arguments for UpdateAds
	args := &query.UpdateAdsArgs{
		ID:     id,
		Status: &status,
	}

	// Update the ad status
	err := query.UpdateAds(ctx.Db, args)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: PostDeactivateAdsHandler update error")
	}

	// Increment metrics for ad deactivation
	collector := metrics.GetCollector()
	if collector != nil {
		collector.IncrementAdDeactivated()
	}

	return gin.H{
		"message": "Ad deactivated successfully",
		"id":      id,
		"status":  status,
	}, http.StatusOK, nil
}
