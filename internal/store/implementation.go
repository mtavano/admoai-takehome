package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

// Store is the database wrapper
type SqlStore struct {
	*sqlx.DB
}

func NewSqlStore(driver, dsn string) (*SqlStore, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	// Maximum Idle Connections
	db.SetMaxIdleConns(20)
	// Idle Connection Timeout
	db.SetConnMaxIdleTime(1 * time.Second)
	// Connection Lifetime
	db.SetConnMaxLifetime(30 * time.Second)

	return &SqlStore{db}, nil
}

func (st *SqlStore) BeginTx(ctx context.Context) (Transactioner, error) {
	tx, err := st.DB.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, errors.Wrap(err, "database: Store.BeginTx st.BeginTxx error")
	}
	return tx, nil
}

//func (st *SqlStore) Exec(query string, params ...interface{}) (sql.Result, error) {
//return st.Exec(query, params...)
//}
