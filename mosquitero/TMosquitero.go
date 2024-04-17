package main

import (
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

// GetMosquitero retorna la única instancia de Mosquitero.
func InitMosquitero(server, username, password string) *Mosquitero {
	mqtonce.Do(func() {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(server)
		opts.SetUsername(username)
		opts.SetPassword(password)
		opts.SetAutoReconnect(true)
		opts.SetKeepAlive(2 * time.Second)
		opts.SetPingTimeout(1 * time.Second)

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

func init() {
	// Ejemplo de uso
	mqttServer := viper.GetString("mosquitero.mqttserver")
	username := viper.GetString("mosquitero.username")
	password := viper.GetString("mosquitero.password")

	mosquitero := InitMosquitero(mqttServer, username, password)

	mosquitero.Send("mi/topic", "mensaje")

}
