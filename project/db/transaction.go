package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func updateInTx(ctx context.Context, db *sqlx.DB, isolation sql.IsolationLevel, fn func(ctx context.Context, tx *sqlx.Tx) error) (err error) {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: isolation,
	})

	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Join(rollbackErr)
			}
			return
		}

		err = tx.Commit()
	}()

	return fn(ctx, tx)
}
