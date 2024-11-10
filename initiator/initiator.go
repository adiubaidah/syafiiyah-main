package initiator

import (
	"context"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/handler"
	"github.com/adiubaidah/rfid-syafiiyah/internal/routing"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/config"
	"github.com/adiubaidah/rfid-syafiiyah/platform/routers"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("santriorder", model.IsValidSantriOrder)
		v.RegisterValidation("userrole", model.IsValidUserRole)
		v.RegisterValidation("userorder", model.IsValidUserOrder)
		v.RegisterValidation("parentorder", model.IsValidParentOrder)
	}

	userUseCase := usecase.NewUserUseCase(store)
	userHandler := handler.NewUserHandler(logger, userUseCase)
	userRouting := routing.UserRouting(userHandler)

	parentUseCase := usecase.NewParentUseCase(store)
	parentHandler := handler.NewParentHandler(&env, logger, parentUseCase, userUseCase)
	parentRouting := routing.ParentRouting(parentHandler)

	santriOccupationUseCase := usecase.NewSantriOccupationUseCase(store)
	santriOccupationHandler := handler.NewSantriOccupationHandler(logger, santriOccupationUseCase)
	santriOccupationRouting := routing.SantriOccupationRouting(santriOccupationHandler)

	santriUseCase := usecase.NewSantriUseCase(store)
	santriHandler := handler.NewSantriHandler(&env, logger, santriUseCase)
	santriRouting := routing.SantriRouting(santriHandler)

	var routerList []routers.Route
	routerList = append(routerList, userRouting...)
	routerList = append(routerList, parentRouting...)
	routerList = append(routerList, santriOccupationRouting...)
	routerList = append(routerList, santriRouting...)

	server := routers.NewRouting(env.ServerAddress, routerList)
	server.Serve()

}
