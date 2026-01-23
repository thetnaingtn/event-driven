package db

import (
	"context"
	"tickets/entities"

	"github.com/jmoiron/sqlx"
)

type DataLake struct {
	db *sqlx.DB
}

func NewDataLake(db *sqlx.DB) DataLake {
	if db == nil {
		panic("db is nil")
	}

	return DataLake{
		db: db,
	}
}

func (d DataLake) Add(ctx context.Context, e entities.Event) error {
	stmt := `
		INSERT INTO 
			events (event_id, published_at, event_name, event_payload)
		VALUES
			(:event_id, :published_at, :event_name, :event_payload)
	`
	_, err := d.db.NamedExecContext(
		ctx,
		stmt,
		e,
	)

	if err != nil {
		return err
	}

	return nil
}
