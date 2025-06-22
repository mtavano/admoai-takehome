package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitSetup, downInitSetup)
}

func upInitSetup(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
		CREATE TABLE ads (
			id TEXT NOT NULL PRIMARY KEY,
			title TEXT NOT NULL,
			image_url TEXT NOT NULL,
			placement TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'active',
			created_at INTEGER NOT NULL,
			ttl integer
		);
	`)

	return err
}

func downInitSetup(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`DROP TABLE ads;`)

	return err
}
