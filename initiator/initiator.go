package initiator

import (
	"context"

	"github.com/adiubaidah/rfid-syafiiyah/internal/handler"
	"github.com/adiubaidah/rfid-syafiiyah/internal/routing"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/config"
	"github.com/adiubaidah/rfid-syafiiyah/platform/routers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func Init() {
	logger := logrus.New()
	env, err := config.LoadConfig(".")
	if err != nil {
		logger.Fatalf("%s cannot load config", err.Error())
	}

	connPool, err := pgxpool.New(context.Background(), env.DBSource)
	if err != nil {
		logger.Fatalf("%s cannot connect to database", err.Error())
	}
	defer connPool.Close()

	store := db.NewStore(connPool)

	santriOccupationUseCase := usecase.NewSantriOccupationUseCase(store)
	santriOccupationHandler := handler.NewSantriOccupationHandler(logger, santriOccupationUseCase)
	santriOccupationRouting := routing.SantriOccupationRouting(santriOccupationHandler)

	var routerList []routers.Route
	routerList = append(routerList, santriOccupationRouting...)

	server := routers.NewRouting(env.ServerAddress, routerList)
	server.Serve()

}
