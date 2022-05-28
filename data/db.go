package data

import (
	"database/sql"
	"errors"
	"log"
	"sync/atomic"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/crypto-pe/api/config"
	"github.com/crypto-pe/api/data/sqlc"
	"github.com/jackc/pgx"
)

var (
	DB        *sqlc.Queries
	connected int32
)

func init() {
	// Set UTC timezone for all our data models & ignore TZ coming from OS env.
	time.Local = time.UTC
}

func PrepareDB(config *config.Config) (err error) {
	tries := 5
	// TODO: Use correct credentials
	conn, err := sql.Open("pgx", config.DBString())
	if err != nil {
		return err
	}

	for tries > 0 {
		log.Println("attempting to make a connection to the database...")
		err = conn.Ping()
		if err != nil {
			tries -= 1
			log.Println(err, "could not connect. retrying...")
			time.Sleep(8 * time.Second)
			continue
		}
		DB = sqlc.New(conn)
		log.Println("connection to the database established.")
		atomic.StoreInt32(&connected, 1)
		return nil
	}
	return errors.New("could not make a connection to the database.")
}

// something
func Disconnect() {
	DB = nil
	atomic.StoreInt32(&connected, 0)
}

func IsConnected() bool {
	return atomic.LoadInt32(&connected) == 1
}

var ErrNoRows = pgx.ErrNoRows
