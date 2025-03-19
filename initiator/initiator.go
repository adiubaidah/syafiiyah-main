package initiator

import (
	"context"
	"time"

	"github.com/adiubaidah/syafiiyah-main/internal/api/handler"
	"github.com/adiubaidah/syafiiyah-main/internal/api/middleware"
	router "github.com/adiubaidah/syafiiyah-main/internal/api/router"
	"github.com/adiubaidah/syafiiyah-main/internal/constant/model"
	mqttHandler "github.com/adiubaidah/syafiiyah-main/internal/mqtt"
	pb "github.com/adiubaidah/syafiiyah-main/internal/protobuf"
	repo "github.com/adiubaidah/syafiiyah-main/internal/repository"
	"github.com/adiubaidah/syafiiyah-main/internal/usecase"
	"github.com/adiubaidah/syafiiyah-main/pkg/config"
	"github.com/adiubaidah/syafiiyah-main/pkg/token"
	"github.com/adiubaidah/syafiiyah-main/platform/mqtt"
	"github.com/adiubaidah/syafiiyah-main/platform/routers"
	storage "github.com/adiubaidah/syafiiyah-main/platform/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	awsCreds "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Init() {

	logger := logrus.New()
	env, err := config.Load(".")
	if err != nil {
		logger.Fatalf("%s cannot load config", err.Error())
	}

	awsConfig, err := awsCfg.LoadDefaultConfig(context.Background(), awsCfg.WithCredentialsProvider(awsCreds.StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID:     env.AWSAccessKey,
			SecretAccessKey: env.AWSSecretKey,
		},
	}), awsCfg.WithRegion(env.AWSRegion))

	awsS3Client := s3.NewFromConfig(awsConfig)
	storageManager := storage.NewStorageManager(awsS3Client, env.AWSBucketName, env.AWSRegion)

	if err != nil {
		logger.Fatalf("Unable to parse awsConfig: %v", err)
	}

	pgConfig, err := pgxpool.ParseConfig(env.DBSource)
	if err != nil {
		logger.Fatalf("Unable to parse pgConfig: %v", err)
	}

	pgConfig.MaxConns = 30
	pgConfig.MinConns = 5
	pgConfig.MaxConnIdleTime = time.Minute * 5
	pgConfig.MaxConnLifetime = time.Hour

	connPool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		logger.Fatalf("Unable to create connection pool: %v", err)
	}
	defer connPool.Close()

	store := repo.NewStore(connPool)
	redisClient := redis.NewClient(&redis.Options{
		DB:   env.DBRedis,
		Addr: env.RedisAddress,
	})
	defer redisClient.Close()

	scheduleServiceConn, err := grpc.NewClient(env.ScheduleServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("Unable to create schedule service connection: %v", err)
	}
	defer scheduleServiceConn.Close()

	if validateActor, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validateActor.RegisterValidation("santri-order", model.IsValidSantriOrder)
		validateActor.RegisterValidation("role", model.IsValidRole)
		validateActor.RegisterValidation("userorder", model.IsValidUserOrder)
		validateActor.RegisterValidation("parentorder", model.IsValidParentOrder)
		validateActor.RegisterValidation("employee-order", model.IsValidEmployeeOrder)
		validateActor.RegisterValidation("valid-time", model.IsValidTime)
		validateActor.RegisterValidation("presencetype", model.IsValidPresenceType)
	}
	tokenMaker, err := token.NewJWTMaker(env.TokenSymmetricKey)
	if err != nil {
		logger.Fatalf("%s cannot create token maker", err.Error())
	}

	sessionUseCase := usecase.NewSessionUseCase(redisClient)

	middle := middleware.NewMiddleware(logger, tokenMaker)

	userUseCase := usecase.NewUserUseCase(store)
	userHandler := handler.NewUserHandler(&handler.UserHandler{
		Logger:  logger,
		UseCase: userUseCase,
	})
	useRouter := router.UserRouter(userHandler)

	authHandler := handler.NewAuthHandler(&handler.AuthHandler{
		Config:         &env,
		UserUseCase:    userUseCase,
		SessionUseCase: sessionUseCase,
		Logger:         logger,
		TokenMaker:     tokenMaker,
	})
	authRouter := router.AuthRouter(middle, authHandler)

	parentUseCase := usecase.NewParentUseCase(store)
	parentHandler := handler.NewParentHandler(&handler.ParentHandler{
		Config:      &env,
		Logger:      logger,
		Storage:     storageManager,
		UseCase:     parentUseCase,
		UserUseCase: userUseCase,
	})
	parentRouter := router.ParentRouter(middle, parentHandler)

	santriScheduleService := pb.NewSantriScheduleServiceClient(scheduleServiceConn)
	santriScheduleHandler := handler.NewSantriScheduleHandler(logger, santriScheduleService)
	santriScheduleRouter := router.SantriScheduleRouter(santriScheduleHandler)

	santriOccupationUseCase := usecase.NewSantriOccupationUseCase(store)
	santriOccupationHandler := handler.NewSantriOccupationHandler(logger, santriOccupationUseCase)
	santriOccupationRouter := router.SantriOccupationRouter(middle, santriOccupationHandler)

	santriUseCase := usecase.NewSantriUseCase(store)
	santriHandler := handler.NewSantriHandler(logger, &env, storageManager, santriUseCase)
	santriRouter := router.SantriRouter(middle, santriHandler)

	santriPresenceUseCase := usecase.NewSantriPresenceUseCase(store)
	santriPresenceHandler := handler.NewSantriPresenceHandler(logger, santriPresenceUseCase)
	santriPresenceRouter := router.SantriPresenceRouter(santriPresenceHandler)

	// employeeScheduleService := pb.NewEmployeeScheduleServiceClient(scheduleServiceConn)
	// employeeScheduleHandler := handler.NewEmployee(logger, employeeScheduleService)

	employeeOccupationUseCase := usecase.NewEmployeeOccupationUseCase(store)
	employeeOccupationHandler := handler.NewEmployeeOccupationHandler(logger, employeeOccupationUseCase)
	employeeOccupationRouter := router.EmployeeOccupationRouter(middle, employeeOccupationHandler)

	employeeUseCase := usecase.NewEmployeeUseCase(store)
	employeeHandler := handler.NewEmployeeHandler(&handler.EmployeeHandler{
		Logger:  logger,
		Storage: storageManager,
		UseCase: employeeUseCase,
	})
	employeeRouter := router.EmployeeRouter(middle, employeeHandler)

	profileHandler := handler.NewProfileHandler(&handler.ProfileHandler{
		Logger:          logger,
		EmployeeUseCase: employeeUseCase,
		ParentUseCase:   parentUseCase,
	})
	profileRouter := router.ProfileRouter(middle, profileHandler)

	smartCardUseCase := usecase.NewSmartCardUseCase(store)
	smartCardHandler := handler.NewSmartCardHandler(logger, smartCardUseCase)
	smartCardRouter := router.SmartCardRouter(smartCardHandler)

	deviceUseCase := usecase.NewDeviceUseCase(store)
	// santriPresenceWorker := worker.NewSantriPresenceWorker(logger, santriPresenceUseCase)

	mqttSantriHandler := mqttHandler.NewSantriMQTTHandler(logger, santriUseCase, santriScheduleService, santriPresenceUseCase)
	// mqttEmployeeHandler := mqttHandler.NewEmployeeMQTTHandler(logger, employeeUseCase, santriScheduleService, santriPresenceUseCase)
	mqttBroker := mqtt.NewMQTTBroker(&mqtt.MQTTBrokerConfig{
		Logger:           logger,
		DeviceUseCase:    deviceUseCase,
		SmartCardUseCase: smartCardUseCase,
		BrokerURL:        env.MQTTBroker,
		SantriHandler:    mqttSantriHandler,
	})
	deviceHandler := handler.NewDeviceHandler(&handler.DeviceHandler{
		Logger:      logger,
		UseCase:     deviceUseCase,
		MqttHandler: mqttBroker,
	})
	deviceRouter := router.DeviceRouter(deviceHandler)

	var routerList []routers.Route
	routerList = append(routerList, authRouter...)
	routerList = append(routerList, useRouter...)
	routerList = append(routerList, parentRouter...)
	routerList = append(routerList, santriScheduleRouter...)
	routerList = append(routerList, santriOccupationRouter...)
	routerList = append(routerList, santriRouter...)
	routerList = append(routerList, santriPresenceRouter...)

	routerList = append(routerList, employeeOccupationRouter...)
	routerList = append(routerList, employeeRouter...)

	routerList = append(routerList, profileRouter...)
	routerList = append(routerList, smartCardRouter...)
	routerList = append(routerList, deviceRouter...)

	server := routers.NewRouting(env.ServerAddress, routerList)
	server.Serve()

}
