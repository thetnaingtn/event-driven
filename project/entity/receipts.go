package entity

type IssueReceiptRequest struct {
	TicketID       string  `json:"ticket_id"`
	Price          Money   `json:"price"`
	IdempotencyKey *string `json:"idempotency_key"`
}

type VoidReceiptRequest struct {
	TicketId       string  `json:"string"`
	IdempotencyKey *string `json:"idempotency_key"`
}
