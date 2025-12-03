package db

import (
	"context"
	"fmt"
	"tickets/entity"

	"github.com/jmoiron/sqlx"
)

type ShowRepository struct {
	db *sqlx.DB
}

func NewShowRepository(db *sqlx.DB) *ShowRepository {
	return &ShowRepository{
		db: db,
	}
}

func (r *ShowRepository) CreateShow(ctx context.Context, show *entity.Show) error {
	stmt := `
		INSERT INTO shows(show_id, dead_nation_id, number_of_tickets, start_time, title, venue)
		VALUES (:show_id, :dead_nation_id, :number_of_tickets, :start_time, :title, :venue) ON CONFLICT DO NOTHING
	`
	_, err := r.db.NamedExecContext(ctx, stmt, show)
	if err != nil {
		return fmt.Errorf("failed to insert show: %w", err)
	}

	return nil
}

func (r *ShowRepository) FindByID(ctx context.Context, id string) (*entity.Show, error) {
	stmt := `
		SELECT show_id, dead_nation_id, number_of_tickets, start_time, title, venue FROM shows
		WHERE show_id=$1;
	`

	res := r.db.QueryRowContext(ctx, stmt, id)
	if res.Err() != nil {
		return nil, fmt.Errorf("failed to query show by id: %w", res.Err())
	}

	var show entity.Show
	if err := res.Scan(&show.ShowID, &show.DeadNationID, &show.NumberOfTickets, &show.StartTime, &show.Title, &show.Venue); err != nil {
		return nil, fmt.Errorf("failed to scan show: %w", err)
	}

	return &show, nil
}
