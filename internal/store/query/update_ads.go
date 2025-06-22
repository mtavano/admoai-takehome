package query

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/mtavano/admoai-takehome/internal/store"
)

type UpdateAdsArgs struct {
	ID        string
	Title     *string
	ImageURL  *string
	Placement *string
	Status    *string
}

func UpdateAds(tx store.Transaction, args *UpdateAdsArgs) error {
	// Validate that ID is provided
	if args.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	// Build update query using squirrel
	query := squirrel.Update("ads").Where(squirrel.Eq{"id": args.ID})

	// Add update fields using COALESCE for conditional updates
	updateMap := make(map[string]interface{})

	updateMap["title"] = squirrel.Expr("COALESCE(?, title)", args.Title)
	updateMap["image_url"] = squirrel.Expr("COALESCE(?, image_url)", args.ImageURL)
	updateMap["placement"] = squirrel.Expr("COALESCE(?, placement)", args.Placement)
	updateMap["status"] = squirrel.Expr("COALESCE(?, status)", args.Status)

	query = query.SetMap(updateMap)

	// Generate SQL query with placeholders for SQLite
	sql, queryArgs, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	// Execute the update
	result, err := tx.Exec(sql, queryArgs...)
	if err != nil {
		return fmt.Errorf("failed to update ads: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no ads found with ID: %s", args.ID)
	}

	return nil
}
