package initiator

import (
	"context"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/handler"
	"github.com/adiubaidah/rfid-syafiiyah/internal/routing"
	"github.com/adiubaidah/rfid-syafiiyah/internal/storage/cache"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/config"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/token"
	"github.com/adiubaidah/rfid-syafiiyah/platform/mqtt"
	"github.com/adiubaidah/rfid-syafiiyah/platform/routers"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func Init() {
	logger := logrus.New()
	env, err := config.LoadConfig(".")
	if err != nil {
		logger.Fatalf("%s cannot load config", err.Error())
	}

	config, err := pgxpool.ParseConfig(env.DBSource)
	if err != nil {
		logger.Fatalf("Unable to parse config: %v", err)
	}

	config.MaxConns = 30
	config.MinConns = 5
	config.MaxConnIdleTime = time.Minute * 5
	config.MaxConnLifetime = time.Hour

	connPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Fatalf("Unable to create connection pool: %v", err)
	}
	defer connPool.Close()

	store := db.NewStore(connPool)
	redisClient := redis.NewClient(&redis.Options{
		DB:   env.DBRedis,
		Addr: env.RedisAddress,
	})
	defer redisClient.Close()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("santriorder", model.IsValidSantriOrder)
		v.RegisterValidation("userrole", model.IsValidUserRole)
		v.RegisterValidation("userorder", model.IsValidUserOrder)
		v.RegisterValidation("parentorder", model.IsValidParentOrder)
		v.RegisterValidation("validTime", model.IsValidTime)
	}
	tokenMaker, err := token.NewJWTMaker(env.TokenSymmetricKey)
	if err != nil {
		logger.Fatalf("%s cannot create token maker", err.Error())
	}

	cacheClient := cache.NewClient(redisClient)

	userUseCase := usecase.NewUserUseCase(store)
	userHandler := handler.NewUserHandler(logger, userUseCase)
	userRouting := routing.UserRouting(userHandler)

	authHandler := handler.NewAuthHandler(userUseCase, cacheClient, &env, logger, tokenMaker)
	authRouting := routing.AuthRouting(authHandler)

	holidayUseCase := usecase.NewHolidayUseCase(store)
	holidayHandler := handler.NewHolidayHandler(logger, holidayUseCase)
	holidayRouting := routing.HolidayRouting(holidayHandler)

	parentUseCase := usecase.NewParentUseCase(store)
	parentHandler := handler.NewParentHandler(&env, logger, parentUseCase, userUseCase)
	parentRouting := routing.ParentRouting(parentHandler)

	santriScheduleUseCase := usecase.NewSantriScheduleUseCase(store)
	santriScheduleHandler := handler.NewSantriScheduleHandler(logger, santriScheduleUseCase)
	santriScheduleRouting := routing.SantriScheduleRouting(santriScheduleHandler)

	santriOccupationUseCase := usecase.NewSantriOccupationUseCase(store)
	santriOccupationHandler := handler.NewSantriOccupationHandler(logger, santriOccupationUseCase)
	santriOccupationRouting := routing.SantriOccupationRouting(santriOccupationHandler)

	santriUseCase := usecase.NewSantriUseCase(store)
	santriHandler := handler.NewSantriHandler(&env, logger, santriUseCase)
	santriRouting := routing.SantriRouting(santriHandler)

	arduinoUseCase := usecase.NewArduinoUseCase(store)
	mqttHandler := mqtt.NewMQTTHandler(arduinoUseCase, env.MQTTBroker)
	arduinoHandler := handler.NewArduinoHandler(logger, arduinoUseCase, mqttHandler)
	arduinoRouting := routing.ArduinoRouting(arduinoHandler)

	var routerList []routers.Route
	routerList = append(routerList, authRouting...)
	routerList = append(routerList, userRouting...)
	routerList = append(routerList, holidayRouting...)
	routerList = append(routerList, parentRouting...)
	routerList = append(routerList, santriScheduleRouting...)
	routerList = append(routerList, santriOccupationRouting...)
	routerList = append(routerList, santriRouting...)
	routerList = append(routerList, arduinoRouting...)

	server := routers.NewRouting(env.ServerAddress, routerList)
	server.Serve()

}
