package db

import (
	"context"
	"fmt"
	"tickets/entity"

	"github.com/jmoiron/sqlx"
)

type TicketRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) *TicketRepository {
	return &TicketRepository{
		db: db,
	}
}

func (r *TicketRepository) SaveTicket(ctx context.Context, ticket *entity.Ticket) (*entity.Ticket, error) {
	if ticket == nil {
		return nil, fmt.Errorf("ticket is nil")
	}

	stmt := `
		INSERT INTO tickets(ticket_id, price_amount, price_currency, customer_email) 
		VALUES (:ticket_id, :price.amount, :price.currency, :customer_email)
	`

	_, err := r.db.NamedExecContext(ctx, stmt, ticket)
	if err != nil {
		return nil, fmt.Errorf("failed to insert ticket: %w", err)
	}

	return ticket, nil
}
