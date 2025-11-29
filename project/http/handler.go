package http

import (
	"context"
	"tickets/entity"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type ticketRepository interface {
	FindAll(ctx context.Context) ([]entity.Ticket, error)
}

type showRepository interface {
	CreateShow(ctx context.Context, show *entity.Show) error
}

type Handler struct {
	eventBus         *cqrs.EventBus
	ticketRepository ticketRepository
	showRepository   showRepository
}
