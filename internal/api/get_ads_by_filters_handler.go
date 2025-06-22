package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mtavano/admoai-takehome/internal/store"
	"github.com/mtavano/admoai-takehome/internal/store/query"
	"github.com/pkg/errors"
)

func GetAdsByFiltersHandler(c *gin.Context, ctx *Context) (any, int, error) {
	// Get query parameters
	placement := c.Query("placement")
	status := c.Query("status")

	var filterExpired bool
	if status == store.AdvertiseStatusInactive {
		filterExpired = true
	}

	// Create arguments for SelectAds
	args := &query.SelectAdsArgs{
		Placement:       placement,
		Status:          status,
		FilterByExpired: filterExpired,
	}

	// Query the database
	records, err := query.SelectAds(ctx.Db, args)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: GetAdsByFiltersHandler query error")
	}

	for idx, rec := range records {
		rec.CalculateAndSetExpired()
		records[idx] = rec
	}

	// Return all records found (could be empty array)
	return map[string]any{
		"ads": records,
	}, http.StatusOK, nil
}
