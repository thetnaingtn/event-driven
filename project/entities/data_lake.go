package entities

import "time"

type DataLakeEvent struct {
	EventID      string    `db:"event_id"`
	EventName    string    `db:"event_name"`
	PublishedAt  time.Time `db:"published_at"`
	EventPayload []byte    `db:"event_payload"`
}
