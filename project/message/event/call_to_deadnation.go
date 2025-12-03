package event

import (
	"context"
	"tickets/entity"
)

func (h Handler) CallToDeadNation(ctx context.Context, event *entity.BookingMade) error {
	show, err := h.showRepository.FindByID(ctx, event.ShowID.String())
	if err != nil {
		return err
	}

	return h.bookingAPIClient.MakeBooking(ctx, entity.CreateBookingRequest{
		BookingID:       event.BookingID,
		EventID:         show.DeadNationID,
		CustomerEmail:   event.CustomerEmail,
		NumberOfTickets: event.NumberOfTickets,
	})
}
