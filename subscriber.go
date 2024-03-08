package main

import (
	"fmt"
	"os"
	"os/signal"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var broker = "tcp://localhost:1883"
var topic = "#"

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Mensagem recebida no tópico: %s\n", message.Topic())
	fmt.Printf("Payload: %s\n", message.Payload())
}

func main() {
	opts := MQTT.NewClientOptions().AddBroker(broker)
	opts.SetClientID("go-mqtt-client")
	opts.SetDefaultPublishHandler(onMessageReceived)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Inscreva-se no tópico
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Printf("Conectado ao broker MQTT em %s\n", broker)

	// Capture sinais para limpeza adequada quando interrompido
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel

	fmt.Println("Desconectando do broker MQTT...")
	client.Disconnect(250)
	fmt.Println("Programa encerrado")
}
