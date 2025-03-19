package mqtt

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/adiubaidah/syafiiyah-main/internal/constant/model"
	mqttHandler "github.com/adiubaidah/syafiiyah-main/internal/mqtt"
	repo "github.com/adiubaidah/syafiiyah-main/internal/repository"
	"github.com/adiubaidah/syafiiyah-main/internal/usecase"
	"github.com/adiubaidah/syafiiyah-main/pkg/util"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type MQTTBroker struct {
	logger           *logrus.Logger
	validator        *validator.Validate
	Client           mqtt.Client
	Topics           map[string]struct{}
	deviceUseCase    *usecase.DeviceUseCase
	smartCardUseCase *usecase.SmartCardUseCase
	SantriHandler    *mqttHandler.SantriMQTTHandler
	mu               sync.Mutex
	MessageHandler   mqtt.MessageHandler
}

type MQTTBrokerConfig struct {
	Logger           *logrus.Logger
	DeviceUseCase    *usecase.DeviceUseCase
	SmartCardUseCase *usecase.SmartCardUseCase
	SantriHandler    *mqttHandler.SantriMQTTHandler
	BrokerURL        string
	IsDevelopment    bool
}

func NewMQTTBroker(config *MQTTBrokerConfig) *MQTTBroker {
	handler := &MQTTBroker{
		logger:           config.Logger,
		validator:        validator.New(),
		Topics:           make(map[string]struct{}),
		deviceUseCase:    config.DeviceUseCase,
		smartCardUseCase: config.SmartCardUseCase,
		SantriHandler:    config.SantriHandler,
	}
	handler.Init(config.BrokerURL)
	handler.RefreshTopics()

	return handler
}

func (h *MQTTBroker) Init(brokerURL string) {
	h.logger.Println("Initializing MQTT client...")
	opts := mqtt.NewClientOptions().AddBroker(brokerURL)
	opts.SetClientID("syafiiyah-main")
	opts.SetDefaultPublishHandler(h.defaultMessageHandler())

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if !token.WaitTimeout(3 * time.Second) {
		h.logger.Fatalf("MQTT broker connection timeout: %v", token.Error())
	}
	if err := token.Error(); err != nil {
		h.logger.Fatalf("Error connecting to MQTT broker: %v", err)
	}

	h.Client = client
	h.logger.Println("Connected to MQTT broker")
}

func (h *MQTTBroker) defaultMessageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		h.logger.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
		if _, exists := h.Topics[msg.Topic()]; !exists {
			h.logger.Warnf("Topic %s is not registered\n", msg.Topic())
			return
		}
		deviceMode := util.GetDeviceMode(msg.Topic())
		deviceName := util.GetDeviceName(msg.Topic())
		acknowledgmentTopic := deviceName + "/acknowledgment/" + deviceMode

		var request model.SmartCardRequest
		err := json.Unmarshal(msg.Payload(), &request)
		if err != nil {
			h.logger.Errorf("Error unmarshaling payload: %v\n", err)

			token := client.Publish(acknowledgmentTopic, 0, false, model.ResponseMessage{
				Code:    400,
				Status:  "error",
				Message: err.Error(),
			})
			if token.Wait() && token.Error() != nil {
				h.logger.Errorf("Error sending acknowledgment: %v\n", token.Error())
			}

			return
		}

		//validate request
		if err := h.validator.Struct(request); err != nil {
			h.logger.Errorf("Error validating request: %v\n", err)
			token := client.Publish(acknowledgmentTopic, 0, false, model.ResponseMessage{
				Code:    400,
				Status:  "error",
				Message: err.Error(),
			})
			if token.Wait() && token.Error() != nil {
				h.logger.Errorf("Error sending acknowledgment: %v\n", token.Error())
			}
			return
		}

		switch repo.DeviceModeType(deviceMode) {
		case repo.DeviceModeTypeRecord:
			h.handleRecord(acknowledgmentTopic, &request)
		case repo.DeviceModeTypePresence:
			h.handlePresence(acknowledgmentTopic, &request)
		case repo.DeviceModeTypePermission:
			h.handlePermission(acknowledgmentTopic, &request)
		case repo.DeviceModeTypePing:
			h.handlePing(acknowledgmentTopic)
		default:
			h.logger.Warnf("Unhandled topic: %s\n", msg.Topic())
		}

	}
}

func (h *MQTTBroker) RefreshTopics() {
	h.logger.Println("Fetching initial topics...")
	device, err := h.deviceUseCase.ListDevices(context.Background())
	if err != nil {
		h.logger.Fatalf("Error fetching device topics: %v", err)
	}

	var newTopics []string
	for _, device := range *device {
		for _, mode := range device.Modes {
			newTopics = append(newTopics, mode.InputTopic)
		}
	}

	// Langganan topic
	h.UpdateSubscriptions(newTopics)
}

func (h *MQTTBroker) UpdateSubscriptions(newTopics []string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// mandekno goroutine untuk topik seng ora diperlukno
	for topic, _ := range h.Topics {
		if !util.Contains(newTopics, topic) {
			h.logger.Printf("Unsubscribing and stopping topic: %s\n", topic)
			h.Client.Unsubscribe(topic)
			delete(h.Topics, topic)
		}
	}

	// Tambah topic anyar
	for _, topic := range newTopics {
		if _, exists := h.Topics[topic]; !exists {
			h.logger.Printf("Subscribing and starting topic: %s\n", topic)
			h.Client.Subscribe(topic, 0, h.defaultMessageHandler())
			h.Topics[topic] = struct{}{}
		}
	}
}
