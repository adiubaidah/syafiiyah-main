package initiator

import (
	"context"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/internal/api/handler"
	"github.com/adiubaidah/rfid-syafiiyah/internal/api/middleware"
	router "github.com/adiubaidah/rfid-syafiiyah/internal/api/router"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	mqttHandler "github.com/adiubaidah/rfid-syafiiyah/internal/mqtt"
	"github.com/adiubaidah/rfid-syafiiyah/internal/storage/cache"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/internal/worker"
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
		validateActor.RegisterValidation("santri-order", model.IsValidSantriOrder)
		validateActor.RegisterValidation("role", model.IsValidRole)
		validateActor.RegisterValidation("userorder", model.IsValidUserOrder)
		validateActor.RegisterValidation("parentorder", model.IsValidParentOrder)
		validateActor.RegisterValidation("employee-order", model.IsValidEmployeeOrder)
		validateActor.RegisterValidation("validTime", model.IsValidTime)
		validateActor.RegisterValidation("presencetype", model.IsValidPresenceType)
	}
	tokenMaker, err := token.NewJWTMaker(env.TokenSymmetricKey)
	if err != nil {
		logger.Fatalf("%s cannot create token maker", err.Error())
	}

	cacheClient := cache.NewClient(redisClient)

	middle := middleware.NewMiddleware(logger, tokenMaker)
	userUseCase := usecase.NewUserUseCase(store)
	userHandler := handler.NewUserHandler(logger, userUseCase)
	useRouter := router.UserRouter(userHandler)

	authHandler := handler.NewAuthHandler(userUseCase, cacheClient, &env, logger, tokenMaker)
	authRouter := router.AuthRouter(middle, authHandler)

	holidayUseCase := usecase.NewHolidayUseCase(store)
	holidayHandler := handler.NewHolidayHandler(logger, holidayUseCase)
	holidayRouter := router.HolidayRouter(middle, holidayHandler)

	parentUseCase := usecase.NewParentUseCase(store)
	parentHandler := handler.NewParentHandler(&env, logger, parentUseCase, userUseCase)
	parentRouter := router.ParentRouter(middle, parentHandler)

	santriScheduleUseCase := usecase.NewSantriScheduleUseCase(store)
	santriScheduleHandler := handler.NewSantriScheduleHandler(logger, santriScheduleUseCase)
	santriScheduleRouter := router.SantriScheduleRouter(santriScheduleHandler)

	santriOccupationUseCase := usecase.NewSantriOccupationUseCase(store)
	santriOccupationHandler := handler.NewSantriOccupationHandler(logger, santriOccupationUseCase)
	santriOccupationRouter := router.SantriOccupationRouter(middle, santriOccupationHandler)

	santriUseCase := usecase.NewSantriUseCase(store)
	santriHandler := handler.NewSantriHandler(logger, &env, santriUseCase)
	santriRouter := router.SantriRouter(middle, santriHandler)

	santriPresenceUseCase := usecase.NewSantriPresenceUseCase(store)
	santriPresenceHandler := handler.NewSantriPresenceHandler(logger, santriPresenceUseCase)
	santriPresenceRouter := router.SantriPresenceRouter(santriPresenceHandler)

	employeeOccupationUseCase := usecase.NewEmployeeOccupationUseCase(store)
	employeeOccupationHandler := handler.NewEmployeeOccupationHandler(logger, employeeOccupationUseCase)
	employeeOccupationRouter := router.EmployeeOccupationRouter(middle, employeeOccupationHandler)

	employeeUseCase := usecase.NewEmployeeUseCase(store)

	profileHandler := handler.NewProfileHandler(logger, employeeUseCase, parentUseCase)
	profileRouter := router.ProfileRouter(middle, profileHandler)

	smartCardUseCase := usecase.NewSmartCardUseCase(store)
	smartCardHandler := handler.NewSmartCardHandler(logger, smartCardUseCase)
	smartCardRouter := router.SmartCardRouter(smartCardHandler)

	deviceUseCase := usecase.NewDeviceUseCase(store)
	santriPresenceWorker := worker.NewSantriPresenceWorker(logger, santriPresenceUseCase)
	scheduleCron := cron.NewScheduleCron(logger, santriScheduleUseCase, santriPresenceWorker)

	mqttSantriHandler := mqttHandler.NewSantriMQTTHandler(logger, scheduleCron, santriUseCase, santriPresenceUseCase)
	mqttBroker := mqtt.NewMQTTBroker(&mqtt.MQTTBrokerConfig{
		Logger:           logger,
		DeviceUseCase:    deviceUseCase,
		SmartCardUseCase: smartCardUseCase,
		BrokerURL:        env.MQTTBroker,
		SantriHandler:    mqttSantriHandler,
	})
	deviceHandler := handler.NewDeviceHandler(logger, deviceUseCase, mqttBroker)
	deviceRouter := router.DeviceRouter(deviceHandler)

	var routerList []routers.Route
	routerList = append(routerList, authRouter...)
	routerList = append(routerList, useRouter...)
	routerList = append(routerList, holidayRouter...)
	routerList = append(routerList, parentRouter...)
	routerList = append(routerList, santriScheduleRouter...)
	routerList = append(routerList, santriOccupationRouter...)
	routerList = append(routerList, santriRouter...)
	routerList = append(routerList, santriPresenceRouter...)
	routerList = append(routerList, employeeOccupationRouter...)

	routerList = append(routerList, profileRouter...)
	routerList = append(routerList, smartCardRouter...)
	routerList = append(routerList, deviceRouter...)

	server := routers.NewRouting(env.ServerAddress, routerList)
	scheduleCron.Start()
	server.Serve()

}
