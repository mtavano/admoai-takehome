package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mtavano/admoai-takehome/internal/store/query"
	"github.com/pkg/errors"
)

func GetAdsByIDHandler(c *gin.Context, ctx *Context) (any, int, error) {
	// Get ID from path parameters
	id := c.Param("id")
	if id == "" {
		return map[string]any{
			"error": "ID parameter is required",
		}, http.StatusBadRequest, nil
	}

	// Create arguments for SelectAds
	args := &query.SelectAdsArgs{
		ID: id,
	}

	// Query the database
	records, err := query.SelectAds(ctx.Db, args)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: GetAdsByIDHandler query error")
	}

	// Check if any records were found
	if len(records) == 0 {
		return map[string]any{
			"error": "Ad not found",
		}, http.StatusNotFound, nil
	}

	// Return the first record (since we're querying by ID, there should be only one)
	return records[0], http.StatusOK, nil
}
