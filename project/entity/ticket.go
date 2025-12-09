package entity

type Ticket struct {
	TicketID      string `json:"ticket_id" db:"ticket_id"`
	Price         Money  `json:"price" db:"price"`
	CustomerEmail string `json:"customer_email" db:"customer_email"`
}

type RefundTicketRequest struct {
	RefundReason   string  `json:"refund_reason"`
	TicketId       string  `json:"string"`
	IdempotencyKey *string `json:"idempotency_key"`
}
