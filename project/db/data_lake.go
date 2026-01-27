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

func (d DataLake) Add(ctx context.Context, event entities.DataLakeEvent) error {
	stmt := `
		INSERT INTO 			    
			events (event_id, published_at, event_name, event_payload) 
		VALUES 
			(:event_id, :published_at, :event_name, :event_payload)`

	_, err := d.db.NamedExecContext(
		ctx,
		stmt,
		event,
	)

	if err != nil {
		return err
	}

	return nil
}

func (d DataLake) GetEvents(ctx context.Context) ([]entities.DataLakeEvent, error) {
	var events []entities.DataLakeEvent
	query := `
		SELECT * FROM events ORDER BY published_at ASC;	
	`

	err := d.db.SelectContext(ctx, &events, query)
	if err != nil {
		return nil, err
	}

	return events, nil
}
