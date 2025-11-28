package db_test

import (
	"log"
	"os"
	"testing"
	"tickets/db"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ticketRepository *db.TicketRepository

func TestMain(m *testing.M) {
	dbConn, err := sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}

	if err := db.InitializeDatabaseSchema(dbConn); err != nil {
		log.Fatal(err)
	}

	ticketRepository = db.NewTicketRepository(dbConn)

	exitVal := m.Run()
	os.Exit(exitVal)
}
