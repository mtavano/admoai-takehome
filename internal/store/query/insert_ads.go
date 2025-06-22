package query

import (
	"github.com/mtavano/admoai-takehome/internal/store"
	"github.com/pkg/errors"
)

func InsertAds(tx store.Transaction, record *store.AdvertiseRecord) error {
	_, err := tx.Exec(`
		INSERT INTO ads (id, title, image_url, placement, status, created_at, expires_at)
		VALUES (?, ?, ?, ?, ?, ?, ?);`,
		record.ID,
		record.Title,
		record.ImageURL,
		record.Placement,
		record.Status,
		record.CreatedAt,
		record.ExpiresAt,
	)
	if err != nil {
		return errors.Wrap(err, "query: InsertAds error")
	}

	return nil
}
