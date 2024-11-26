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
	"github.com/adiubaidah/rfid-syafiiyah/platform/cron"
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

	if validateActor, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validateActor.RegisterValidation("santriorder", model.IsValidSantriOrder)
		validateActor.RegisterValidation("role", model.IsValidRole)
		validateActor.RegisterValidation("userorder", model.IsValidUserOrder)
		validateActor.RegisterValidation("parentorder", model.IsValidParentOrder)
		validateActor.RegisterValidation("validTime", model.IsValidTime)
		validateActor.RegisterValidation("presencetype", model.IsValidPresenceType)
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

	santriPresenceUseCase := usecase.NewSantriPresenceUseCase(store)
	santriPresenceHandler := handler.NewSantriPresenceHandler(logger, santriPresenceUseCase)
	santriPresenceRouting := routing.SantriPresenceRouting(santriPresenceHandler)

	smartCardUseCase := usecase.NewSmartCardUseCase(store)
	smartCardHandler := handler.NewSmartCardHandler(logger, smartCardUseCase)
	smartCardRouting := routing.SmartCardRouting(smartCardHandler)

	deviceUseCase := usecase.NewDeviceUseCase(store)
	scheduleCron := cron.NewScheduleCron(logger, santriScheduleUseCase)
	mqttHandler := mqtt.NewMQTTHandler(&mqtt.MQTTHandlerConfig{
		Logger:                logger,
		DeviceUseCase:         deviceUseCase,
		Schedule:              scheduleCron,
		SmartCardUseCase:      smartCardUseCase,
		SantriUseCase:         santriUseCase,
		SantriPresenceUseCase: santriPresenceUseCase,
		BrokerURL:             env.MQTTBroker,
		// IsDevelopment: env.Gi,
	})
	deviceHandler := handler.NewDeviceHandler(logger, deviceUseCase, mqttHandler)
	deviceRouting := routing.DeviceRouting(deviceHandler)

	var routerList []routers.Route
	routerList = append(routerList, authRouting...)
	routerList = append(routerList, userRouting...)
	routerList = append(routerList, holidayRouting...)
	routerList = append(routerList, parentRouting...)
	routerList = append(routerList, santriScheduleRouting...)
	routerList = append(routerList, santriOccupationRouting...)
	routerList = append(routerList, santriRouting...)
	routerList = append(routerList, santriPresenceRouting...)
	routerList = append(routerList, smartCardRouting...)
	routerList = append(routerList, deviceRouting...)

	server := routers.NewRouting(env.ServerAddress, routerList)
	scheduleCron.Start()
	server.Serve()

}
