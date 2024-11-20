package mqtt

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTHandler struct {
	Client         mqtt.Client
	Topics         map[string]chan struct{}
	Usecase        usecase.DeviceUseCase
	mu             sync.Mutex
	MessageHandler mqtt.MessageHandler
}

func NewMQTTHandler(usecase usecase.DeviceUseCase, brokerURL string) *MQTTHandler {
	handler := &MQTTHandler{
		Topics:  make(map[string]chan struct{}),
		Usecase: usecase,
	}

	handler.Init(brokerURL)

	handler.RefreshTopics()

	return handler
}

func (h *MQTTHandler) Init(brokerURL string) {
	log.Println("Initializing MQTT client...")
	opts := mqtt.NewClientOptions().AddBroker(brokerURL)
	opts.SetClientID("arduino-server")
	opts.SetDefaultPublishHandler(h.defaultMessageHandler())

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if !token.WaitTimeout(3 * time.Second) {
		log.Fatalf("MQTT broker connection timeout: %v", token.Error())
	}
	if err := token.Error(); err != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", err)
	}

	h.Client = client
	log.Println("Connected to MQTT broker")
}

func (h *MQTTHandler) defaultMessageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}
}

func (h *MQTTHandler) RefreshTopics() {
	log.Println("Fetching initial topics...")
	arduinos, err := h.Usecase.ListDevices(context.Background())
	if err != nil {
		log.Fatalf("Error fetching arduino topics: %v", err)
	}

	var newTopics []string
	for _, arduino := range arduinos {
		for _, mode := range arduino.Modes {
			newTopics = append(newTopics, mode.InputTopic)
		}
	}

	// Langganan topic
	h.UpdateSubscriptions(newTopics)
}

func (h *MQTTHandler) UpdateSubscriptions(newTopics []string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// mandekno goroutine untuk topik seng ora diperlukno
	for topic, stopChan := range h.Topics {
		if !util.Contains(newTopics, topic) {
			log.Printf("Unsubscribing and stopping topic: %s\n", topic)
			h.Client.Unsubscribe(topic)
			close(stopChan)
			delete(h.Topics, topic)
		}
	}

	// Tambah topic anyar
	for _, topic := range newTopics {
		if _, exists := h.Topics[topic]; !exists {
			log.Printf("Subscribing and starting topic: %s\n", topic)
			h.Client.Subscribe(topic, 0, h.defaultMessageHandler())
			stopChan := make(chan struct{})
			h.Topics[topic] = stopChan
			go h.listenTopic(topic, stopChan)
		}
	}
}

func (h *MQTTHandler) listenTopic(topic string, stopChan chan struct{}) {
	log.Printf("Listening on topic: %s\n", topic)
	for {
		select {
		case <-stopChan:
			log.Printf("Stopping listener for topic: %s\n", topic)
			return
		}
	}
}
