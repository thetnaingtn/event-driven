package event

import (
	"context"
	"fmt"
	"tickets/entity"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
)

func (h *Handler) PrintTicket(ctx context.Context, event *entity.TicketBookingConfirmed) error {
	log.FromContext(ctx).Info("Printing ticket")

	ticketHTML := `
		<html>
			<head>
				<title>Ticket</title>
			</head>
			<body>
				<h1>Ticket ` + event.TicketID + `</h1>
				<p>Price: ` + event.Price.Amount + ` ` + event.Price.Currency + `</p>	
			</body>
		</html>
`

	err := h.fileAPIClient.UploadFile(ctx, fmt.Sprintf("%s-ticket.html", event.TicketID), ticketHTML)
	if err != nil {
		return fmt.Errorf("failed to upload ticket file: %w", err)
	}

	return nil
}
