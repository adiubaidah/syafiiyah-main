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
	Topics         map[string]struct{}
	UpdateChan     chan struct{}
	Usecase        usecase.ArduinoUseCase
	mu             sync.Mutex
	MessageHandler mqtt.MessageHandler
}

func NewMQTTHandler(usecase usecase.ArduinoUseCase, brokerURL string) *MQTTHandler {
	handler := &MQTTHandler{
		Topics:     make(map[string]struct{}),
		UpdateChan: make(chan struct{}),
		Usecase:    usecase,
	}

	//
	go handler.RunListener()

	handler.Init(brokerURL)

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

	// Trigger pembaruan awal
	log.Println("Sending initial update signal...")
	h.UpdateChan <- struct{}{}
	log.Println("Initial update signal sent.")
}

func (h *MQTTHandler) defaultMessageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}
}

func (h *MQTTHandler) UpdateSubscriptions(newTopics []string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Gusek topic seng ora di-subscribe
	for topic := range h.Topics {
		if !util.Contains(newTopics, topic) {
			h.Client.Unsubscribe(topic)
			delete(h.Topics, topic)
		}
	}

	// Tambah topic anyar
	for _, topic := range newTopics {
		if _, exists := h.Topics[topic]; !exists {
			log.Println("Subscribing to topic", topic)
			h.Client.Subscribe(topic, 0, h.defaultMessageHandler())
			h.Topics[topic] = struct{}{}
		}
	}
}

func (h *MQTTHandler) RunListener() {
	log.Println("Starting MQTT Listener...")
	for {
		select {
		case <-h.UpdateChan:
			log.Println("Update signal received. Fetching new topics...")
			arduinos, err := h.Usecase.ListArduinos(context.Background())
			if err != nil {
				log.Printf("Error fetching arduino topics: %v\n", err)
				continue
			}

			var newTopics []string
			for _, arduino := range arduinos {
				for _, mode := range arduino.Modes {
					newTopics = append(newTopics, mode.InputTopic)
				}
			}
			h.UpdateSubscriptions(newTopics)
		}
	}
}
