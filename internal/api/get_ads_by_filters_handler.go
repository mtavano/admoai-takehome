package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mtavano/admoai-takehome/internal/store/query"
	"github.com/pkg/errors"
)

func GetAdsByFiltersHandler(c *gin.Context, ctx *Context) (any, int, error) {
	// Get query parameters
	placement := c.Query("placement")
	status := c.Query("status")

	// Create arguments for SelectAds
	args := &query.SelectAdsArgs{
		Placement: placement,
		Status:    status,
	}

	// Query the database
	records, err := query.SelectAds(ctx.Db, args)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: GetAdsByFiltersHandler query error")
	}

	// Return all records found (could be empty array)
	return gin.H{
		"ads": records,
		"count": len(records),
	}, http.StatusOK, nil
}
