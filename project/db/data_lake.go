package db

import (
	"context"
	"time"

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

func (d DataLake) Add(ctx context.Context, eventName, eventId string, payload []byte, publishedAt time.Time) error {
	stmt := `
		INSERT INTO 
			events (event_id, published_at, event_name, event_payload)
		VALUES
			($1, $2, $3, $4)
	`
	_, err := d.db.ExecContext(
		ctx,
		stmt,
		eventId,
		publishedAt,
		eventName,
		payload,
	)

	if err != nil {
		return err
	}

	return nil
}
