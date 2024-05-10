package mosquitero

import (
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

type Mosquitero struct {
	client mqtt.Client
}

var mqtinstance *Mosquitero
var mqtonce sync.Once

// InitMosquitero inicializa la única instancia de Mosquitero.
func InitMosquitero(server, username, password string) *Mosquitero {
	mqtonce.Do(func() {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(server)
		opts.SetClientID("go_mqtt_client")
		opts.SetUsername(username)
		opts.SetPassword(password)
		opts.SetAutoReconnect(true)
		opts.SetKeepAlive(2 * time.Second)
		opts.SetPingTimeout(1 * time.Second)
		opts.SetConnectTimeout(5 * time.Second) // Tiempo de espera para la conexión inicial
		opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
			fmt.Printf("Connection lost: %v. Reconnecting...\n", err)

		})

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}

		mqtinstance = &Mosquitero{client: client}
	})
	return mqtinstance
}

func GetMosquitero() *Mosquitero {
	return mqtinstance
}

// Send envía un mensaje al topic especificado.
func (m *Mosquitero) Send(topic string, payload string) {
	go m.InternalSend(topic, payload)
}

func (m *Mosquitero) GetClient() mqtt.Client {
	return m.client
}

func (m *Mosquitero) InternalSend(topic string, payload string) {
	token := m.client.Publish(topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		log.Println("Error al publicar:", token.Error())
	}
}

// CheckConnection verifica y reconecta si es necesario.
func (m *Mosquitero) CheckConnection() {
	if !m.client.IsConnected() {
		if token := m.client.Connect(); token.Wait() && token.Error() != nil {
			log.Println("Error al reconectar:", token.Error())
		}
	}
}

// Subscribe se suscribe a una lista de topics con un handler para los mensajes.
func (m *Mosquitero) Subscribe(topics []string, handler mqtt.MessageHandler) {
	for _, topic := range topics {
		if token := m.client.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
			log.Printf("Error al suscribirse al topic %s: %s", topic, token.Error())
		} else {
			log.Printf("Suscrito al topic %s", topic)
		}
	}
}

func Mosquiteroinit() {
	// Ejemplo de uso
	mqttServer := viper.GetString("mosquitero.mqttserver")
	username := viper.GetString("mosquitero.username")
	password := viper.GetString("mosquitero.password")

	mosquitero := InitMosquitero(mqttServer, username, password)

	// Suscribirse a topics
	topics := []string{"mi/topic", "otro/topic"}
	mosquitero.Subscribe(topics, defaultMessageHandler)

	// Enviar un mensaje
	mosquitero.Send("mi/topic", "mensaje")
}

func defaultMessageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}
