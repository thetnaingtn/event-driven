package entity

type RefundTicket struct {
	Header   MessageHeader `json:"header"`
	TicketID string        `json:"ticket_id"`
}
