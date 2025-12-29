package db_test

import (
	"context"
	"testing"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddTicket(t *testing.T) {
	ticketId := uuid.New()

	ticket := entities.Ticket{
		TicketID: ticketId.String(),
		Price: entities.Money{
			Amount:   "2.5",
			Currency: "USD",
		},
	}

	ctx := context.Background()

	err := ticketRepository.Add(ctx, ticket)

	assert.Nil(t, err)

	tickets, err := ticketRepository.FindAll(ctx)

	assert.Nil(t, err)
	assert.Len(t, tickets, 1)

	// create ticket with same uuid
	err = ticketRepository.Add(ctx, ticket)
	assert.Nil(t, err) // shouldn't throw any error just ignore

	tickets, err = ticketRepository.FindAll(ctx)
	assert.Nil(t, err)
	assert.Len(t, tickets, 1)
}
