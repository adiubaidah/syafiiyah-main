package mqtt

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var UpdateChannel = make(chan struct{})
var Topics map[string]struct{}
var Client mqtt.Client
var Mu sync.Mutex

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

func Init(usecase usecase.ArduinoUseCase, brokerURL string) {
	Topics = make(map[string]struct{})
	opts := mqtt.NewClientOptions().AddBroker(brokerURL)
	opts.SetClientID("arduino-server")
	opts.SetDefaultPublishHandler(messagePubHandler)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Assign ke variabel global
	Client = client

	log.Default().Println("Connected to MQTT broker")

	UpdateChannel <- struct{}{}

}

func UpdateSubscriptions(newTopics []string) {
	Mu.Lock()
	defer Mu.Unlock()

	if Client == nil {
		log.Default().Println("MQTT client is not initialized")
		return
	}

	// Hapus topik yang tidak lagi relevan
	for topic := range Topics {
		if !util.Contains(newTopics, topic) {
			log.Default().Println("Unsubscribing from topic", topic)
			Client.Unsubscribe(topic)
			delete(Topics, topic)
		}
	}

	// Tambahkan topik baru yang belum di-subscribe
	for _, topic := range newTopics {
		if _, exists := Topics[topic]; !exists {
			log.Default().Println("Subscribing to topic", topic)
			Client.Subscribe(topic, 0, messagePubHandler)
			Topics[topic] = struct{}{}
		}
	}
}

func RunMQTTListener(usecase usecase.ArduinoUseCase) {
	for {
		// Tunggu sinyal pembaruan
		<-UpdateChannel

		// Ambil semua topik terbaru dari database
		arduinos, err := usecase.ListArduinos(context.Background())

		var newTopics []string

		for _, arduino := range arduinos {
			for _, mode := range arduino.Modes {
				newTopics = append(newTopics, mode.InputTopic)
			}
		}
		log.Default().Println("New topics", newTopics)

		if err != nil {
			// Log jika terjadi error
			log.Default().Println(err)
			continue
		}

		// Update topik yang didengarkan
		UpdateSubscriptions(newTopics)
	}
}
