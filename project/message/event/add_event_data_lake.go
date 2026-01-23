package event

import (
	"encoding/json"
	"tickets/entities"

	"github.com/ThreeDotsLabs/watermill/message"
)

type Event struct {
	Header entities.MessageHeader `json:"header"`
}

func (h Handler) AddEventToDataLake(msg *message.Message) error {
	var event Event

	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	return h.dataLakeRepository.Add(msg.Context(), marshaler.NameFromMessage(msg), event.Header.ID, msg.Payload, event.Header.PublishedAt)
}
