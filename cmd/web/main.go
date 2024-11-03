package main

import (
	"context"

	sqlcDB "github.com/adiubaidah/rfid-syafiiyah/db/sqlc"
	"github.com/adiubaidah/rfid-syafiiyah/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// db "github.com/adiubaidah/rfid-syafiiyah/db/sqlc"

func main() {
	viperEnv, err := config.LoadEnv("../../..")
	log := config.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	connPool, err := pgxpool.New(context.Background(), viperEnv.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	db := sqlcDB.New(connPool)

}
