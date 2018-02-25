package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/brian-armstrong/gpio"
	"github.com/charles-d-burton/aws-mqtt/messages"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type mqttPlugin struct {
	PubTopic func(topic string, qos byte, retained bool, payload interface{}) MQTT.Token
}

type payload struct {
	Message string `json:"message"`
}

func (plugin mqttPlugin) PluginID() string {
	return "nightmare-trigger"
}

func (plugin mqttPlugin) Topic() string {
	return "/nightmare-cat/button_pressed"
}

func (plugin mqttPlugin) ProcessMessage(msg MQTT.Message) error {
	return nil
}

//TODO:  Need to make plugins configurable
//PublishTopic listen for changes from the hall effect sensor, then publish to MQTT
func (plugin mqttPlugin) PublishTopic(f func(string, byte, bool, interface{}) MQTT.Token) error {
	plugin.PubTopic = f //Set the struct value so we can take action on it
	watcher := gpio.NewWatcher()
	watcher.AddPin(14) //Watch the pin the hall effect sensor is on
	defer watcher.Close()
	for {
		pin, value := watcher.Watch()
		if value == 0 {
			var message payload
			message.Message = "doorbell-pressed"
			b, err := json.Marshal(message)
			if err != nil {
				return fmt.Errorf("Unable to marshal struct")
			}
			if token := f(plugin.Topic(), 0, false, string(b)); token.Wait() && token.Error() != nil {
				log.Println(token.Error())
			}
			log.Println("Published messaged!")
		}
		fmt.Printf("read %d from gpio %d\n", value, pin)
	}
}

func GetPlugin() (messages.MessageReceiver, error) {
	receiver := mqttPlugin{}
	return receiver, nil
}
