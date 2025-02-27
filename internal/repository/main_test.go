package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store
var sqlStore *SQLStore

func TestMain(m *testing.M) {
	env, err := config.Load("../..")
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	connPool, err := pgxpool.New(context.Background(), env.DBSource)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	defer connPool.Close()

	testStore = NewStore(connPool)

	var ok bool
	sqlStore, ok = testStore.(*SQLStore)
	if !ok {
		log.Fatal("Cannot convert to *SQLStore")
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}
