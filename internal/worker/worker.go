package worker

import (
	"context"

	pb "github.com/adiubaidah/rfid-syafiiyah/internal/protobuf"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type EmployeeWorker interface {
	SubscribeSchedule(ctx context.Context)
}

type worker struct {
	client          *redis.Client
	logger          *logrus.Logger
	schedule        pb.EmployeeScheduleServiceClient
	presenceUseCase usecase.EmployeePresenceUseCase
}

func NewEmployeeWorker(logger *logrus.Logger, redisClient *redis.Client, schedule pb.EmployeeScheduleServiceClient, usecase usecase.EmployeePresenceUseCase) EmployeeWorker {
	w := &worker{
		logger:          logger,
		client:          redisClient,
		schedule:        schedule,
		presenceUseCase: usecase,
	}

	go w.SubscribeSchedule(context.Background())

	return w
}

func (w *worker) SubscribeSchedule(ctx context.Context) {
	sub := w.client.Subscribe(ctx, "employee_schedule")
	ch := sub.Channel()

	for msg := range ch {
		w.logger.Info("Message received: ", msg.Payload)
	}
}
