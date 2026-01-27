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

	return h.dataLakeRepository.Add(msg.Context(), entities.DataLakeEvent{
		EventID:      event.Header.ID,
		EventName:    marshaler.NameFromMessage(msg),
		EventPayload: msg.Payload,
		PublishedAt:  event.Header.PublishedAt,
	})
}
