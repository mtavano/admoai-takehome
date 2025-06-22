package query

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/mtavano/admoai-takehome/internal/store"
)

type SelectAdsArgs struct {
	ID        string
	Title     string
	Status    string
	Placement string
}

func SelectAds(tx store.Transaction, args *SelectAdsArgs) ([]*store.AdvertiseRecord, error) {
	// Build query using squirrel
	query := squirrel.Select("*").From("ads")

	// Add conditions based on provided fields
	if args.ID != "" {
		query = query.Where(squirrel.Eq{"id": args.ID})
	}
	if args.Title != "" {
		query = query.Where(squirrel.Eq{"title": args.Title})
	}
	if args.Status != "" {
		query = query.Where(squirrel.Eq{"status": args.Status})
	}
	if args.Placement != "" {
		query = query.Where(squirrel.Eq{"placement": args.Placement})
	}

	// Generate SQL query with placeholders for SQLite
	sql, queryArgs, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	// Execute the query
	record := make([]*store.AdvertiseRecord, 0)
	err = tx.Select(&record, sql, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to select ads: %w", err)
	}

	return record, nil
}
