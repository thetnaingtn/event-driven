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
