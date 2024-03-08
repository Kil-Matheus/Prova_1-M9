package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Configuração do cliente MQTT
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("Loja_1")
	client := MQTT.NewClient(opts)

	// Conecte ao broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)

	// Captura de sinal para lidar com o encerramento
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	// Loop para publicar mensagens continuamente
	go func() {
		for {
			temp, equipamento := simularLeituraSensorTemperaturaLJ1()

			if temp > -25 && temp < -15 {
				message := fmt.Sprintf("Lj 1: Freezer: %.d | Temperatura: %.2f", equipamento, temp)
				token := client.Publish("id = id: lj01f01/freezer", 0, false, message)
				token.Wait()
				fmt.Printf("Publicado: %s\n", message)
				time.Sleep(2 * time.Second)
			}
			if temp > -15 {
				message := fmt.Sprintf("Lj 1: Freezer: %.d | Temperatura: %.2f [ALERTA: Temperatura ALTA]", equipamento, temp)
				token := client.Publish("id = id: lj01f01/freezer", 0, false, message)
				token.Wait()
				fmt.Printf("Publicado: %s\n", message)
				time.Sleep(2 * time.Second)
			}
			if temp < -25 {
				message := fmt.Sprintf("Lj 1: Freezer: %.d | Temperatura: %.2f [ALERTA: Temperatura BAIXA]", equipamento, temp)
				token := client.Publish("id = id: lj01f01/freezer", 0, false, message)
				token.Wait()
				fmt.Printf("Publicado: %s\n", message)
				time.Sleep(2 * time.Second)
			}
		}
	}()

	// Aguarde um sinal de interrupção
	<-signalChannel
	fmt.Println("Publicação encerrada")
}

// simularLeituraSensor simula a leitura do sensor
func simularLeituraSensorTemperaturaLJ1() (float64, int) {
	// Simulação de leitura de sensor
	temp := rand.Float64() * -50
	equipamento := rand.Intn(5)
	return temp, equipamento
}
