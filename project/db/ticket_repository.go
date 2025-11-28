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
		VALUES (:ticket_id, :price.amount, :price.currency, :customer_email) ON CONFLICT DO NOTHING
	`

	_, err := r.db.NamedExecContext(ctx, stmt, ticket)
	if err != nil {
		return nil, fmt.Errorf("failed to insert ticket: %w", err)
	}

	return ticket, nil
}

func (r *TicketRepository) RemoveTicket(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is empty")
	}

	stmt := `
		DELETE FROM tickets WHERE ticket_id=$1;
	`

	_, err := r.db.ExecContext(ctx, stmt, id)

	if err != nil {
		return err
	}

	return nil
}

func (r *TicketRepository) FindAll(ctx context.Context) ([]entity.Ticket, error) {
	stmt := `
		SELECT 
			ticket_id,
			price_amount as "price.amount",
			price_currency as "price.currency",
			customer_email 
		FROM tickets
	`

	var tickets []entity.Ticket

	err := r.db.SelectContext(ctx, &tickets, stmt)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}
