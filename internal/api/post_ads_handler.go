package api

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mtavano/admoai-takehome/internal/store"
	"github.com/mtavano/admoai-takehome/internal/store/query"
	"github.com/pkg/errors"
)

type PostAdsHandlerRequest struct {
	Title     string `json:"title" binding:"required"`
	ImageURL  string `json:"image_url" binding:"required,url"`
	Placement string `json:"placement" binding:"required"`
	Ttl       *uint  `json:"ttl"`
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

	rec := &store.AdvertiseRecord{
		ID:        uuid.NewString(),
		Title:     req.Title,
		ImageURL:  req.ImageURL,
		Placement: req.Placement,
		Status:    store.AdvertiseStatusActive,
		CreatedAt: int64(time.Now().Unix()),
	}

	err := query.InsertAds(ctx.Db, rec)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: PostAdsHandlerRequest error")
	}

	return rec, http.StatusCreated, nil
}
