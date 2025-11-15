package main

import (
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type ProductOutOfStock struct {
	MessageHeader MessageHeader `json:"header"`
	ProductID     string        `json:"product_id"`
}

type ProductBackInStock struct {
	MessageHeader MessageHeader `json:"header"`
	ProductID     string        `json:"product_id"`
	Quantity      int           `json:"quantity"`
}

type MessageHeader struct {
	ID         string `json:"id"`
	EventName  string `json:"event_name"`
	OccurredAt string `json:"occurred_at"`
}

func NewMessageHeader(id, eventName string) MessageHeader {
	return MessageHeader{
		ID:         id,
		EventName:  eventName,
		OccurredAt: time.Now().Format(time.RFC3339),
	}
}

type Publisher struct {
	pub message.Publisher
}

func NewPublisher(pub message.Publisher) Publisher {
	return Publisher{
		pub: pub,
	}
}

func (p Publisher) PublishProductOutOfStock(productID string) error {
	uuid := watermill.NewUUID()

	event := ProductOutOfStock{
		MessageHeader: NewMessageHeader(uuid, "ProductOutOfStock"),
		ProductID:     productID,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := message.NewMessage(uuid, payload)

	return p.pub.Publish("product-updates", msg)
}

func (p Publisher) PublishProductBackInStock(productID string, quantity int) error {
	uuid := watermill.NewUUID()

	event := ProductBackInStock{
		MessageHeader: NewMessageHeader(uuid, "ProductBackInStock"),
		ProductID:     productID,
		Quantity:      quantity,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := message.NewMessage(uuid, payload)

	return p.pub.Publish("product-updates", msg)
}
