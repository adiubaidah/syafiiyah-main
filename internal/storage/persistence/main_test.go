package persistence

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
	// Load configuration
	env, err := config.LoadConfig("../../..")
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	// Initialize connection pool
	connPool, err := pgxpool.New(context.Background(), env.DBSource)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	defer connPool.Close()

	// Initialize store
	testStore = NewStore(connPool)

	// Type assertion to *SQLStore
	var ok bool
	sqlStore, ok = testStore.(*SQLStore)
	if !ok {
		log.Fatal("Cannot convert to *SQLStore")
	}

	// Run tests
	exitCode := m.Run()
	os.Exit(exitCode)
}
