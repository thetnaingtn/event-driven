package outbox

import (
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/jmoiron/sqlx"
)

func NewPostgresSubscriber(db *sqlx.DB, logger watermill.LoggerAdapter) *watermillSQL.Subscriber {
	subscriber, err := watermillSQL.NewSubscriber(db, watermillSQL.SubscriberConfig{
		PollInterval:     time.Millisecond * 100,
		SchemaAdapter:    watermillSQL.DefaultPostgreSQLSchema{},
		OffsetsAdapter:   watermillSQL.DefaultPostgreSQLOffsetsAdapter{},
		InitializeSchema: true,
	}, logger)

	if err != nil {
		panic(fmt.Errorf("failed to create new watermill sql subscriber: %w", err))
	}

	return subscriber
}
