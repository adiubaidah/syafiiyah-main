package bootstrap

import (
	"context"

	db "github.com/adiubaidah/rfid-syafiiyah/db/sqlc"
	"github.com/adiubaidah/rfid-syafiiyah/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Bootstrap struct {
	Queries *db.Queries
	Log     *logrus.Logger
}

func NewBootstrap(path string) (*Bootstrap, error) {
	env, err := config.LoadEnv(path)
	if err != nil {
		return nil, err
	}
	logger := config.NewLogger()
	connPool, err := pgxpool.New(context.Background(), env.DBSource)
	if err != nil {
		return nil, err
	}
	database := db.New(connPool)

	return &Bootstrap{
		Queries: database,
		Log:     logger,
	}, err
}
