package main

import (
	"atividade-5/base"
	"fmt"
	"log"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	broker       = "tcp://localhost:1883"
	requestTopic = "file/request"
	serverDir    = "../files_server"
)

type ServerMQTT struct {
	*base.BaseServer
}

func handleRequest(client mqtt.Client, msg mqtt.Message) {
	server := &ServerMQTT{&base.BaseServer{}}
	payload := string(msg.Payload())
	parts := strings.SplitN(payload, " ", 4)

	if len(parts) < 2 {
		return // comando inválido
	}

	clientID := parts[0]
	command := parts[1]

	filename := ""
	if len(parts) >= 3 {
		filename = parts[2]
	}

	data := ""
	if len(parts) == 4 {
		data = parts[3]
	}

	response := server.HandleRequest(command, filename, data)

	client.Publish("file/response/"+clientID, 0, false, response)
}

func main() {
	os.MkdirAll(serverDir, os.ModePerm)

	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("file-server")
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro ao conectar: %v", token.Error())
		os.Exit(1)
	}

	if token := client.Subscribe(requestTopic, 0, handleRequest); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro ao se inscrever n topico: %v", token.Error())
	}

	fmt.Println("MQTT server is running")
	select {}
}
