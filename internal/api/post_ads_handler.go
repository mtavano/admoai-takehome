package api

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mtavano/admoai-takehome/internal/metrics"
	"github.com/mtavano/admoai-takehome/internal/store"
	"github.com/mtavano/admoai-takehome/internal/store/query"
	"github.com/pkg/errors"
)

type PostAdsHandlerRequest struct {
	Title     string `json:"title" binding:"required"`
	ImageURL  string `json:"image_url" binding:"required,url"`
	Placement string `json:"placement" binding:"required"`
	Ttl       int64  `json:"ttl"`
}

func PostAdsHandler(c *gin.Context, ctx *Context) (any, int, error) {
	var req PostAdsHandlerRequest

	// Bind JSON with validation
	if err := c.ShouldBindJSON(&req); err != nil {
		return map[string]any{
			"error":   "Validation failed",
			"details": err.Error(),
		}, http.StatusBadRequest, nil
	}

	// Additional custom validation for image_url
	if _, err := url.ParseRequestURI(req.ImageURL); err != nil {
		return map[string]any{
			"error":   "Invalid image URL format",
			"details": "The image_url must be a valid URL",
		}, http.StatusBadRequest, nil
	}

	createdAt := time.Now()
	var expiresAt *int64
	if req.Ttl > 0 {
		// Calculate expiration time by adding TTL minutes to creation time
		expAt := createdAt.Add(time.Duration(req.Ttl) * time.Minute)
		expAtUnix := int64(expAt.Unix())
		expiresAt = &expAtUnix
	}

	rec := &store.AdvertiseRecord{
		ID:        uuid.NewString(),
		Title:     req.Title,
		ImageURL:  req.ImageURL,
		Placement: req.Placement,
		Status:    store.AdvertiseStatusActive,
		CreatedAt: int64(createdAt.Unix()),
		ExpiresAt: expiresAt,
	}

	err := query.InsertAds(ctx.Db, rec)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: PostAdsHandlerRequest error")
	}

	// Increment metrics for ad creation
	collector := metrics.GetCollector()
	if collector != nil {
		collector.IncrementAdCreated()
	}

	return rec, http.StatusCreated, nil
}
