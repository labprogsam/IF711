package main

import (
	"atividade-5/base"
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	broker       = "tcp://localhost:1883"
	requestTopic = "file/request"
	clientFiles  = "../files_client"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type ClientMQTT struct {
	base *base.BaseClient
}

func NewClientRabbitMQ() *ClientMQTT {
	c := &ClientMQTT{}
	c.base = base.NewBaseClient(c)
	return c
}

type MQTTParams struct {
	client   mqtt.Client
	clientID string
}

var responseChan = make(chan string)

func (c *ClientMQTT) SendCommand(command string, input any) string {
	params, ok := input.(MQTTParams)
	if !ok {
		fmt.Println("parametros errados")
	}
	payload := fmt.Sprintf("%s %s", params.clientID, command)
	params.client.Publish(requestTopic, 0, false, payload)

	for r := range responseChan {
		return r
	}
	return ""
}

func (c *ClientMQTT) HandleCommand(cmd []string, params MQTTParams) {
	c.base.HandleCommand(cmd, params)
}

func main() {
	os.MkdirAll(clientFiles, os.ModePerm)
	mClient := NewClientRabbitMQ()

	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Usage: <LIST | UPLOAD | DOWNLOAD> <ID> <FILENAME>")
		return
	}

	clientId := args[1]

	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientId)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		failOnError(token.Error(), "Failed to connect to MQTT")
	}

	// Subscribe to personal response topic
	respTopic := "file/response/" + clientId
	client.Subscribe(respTopic, 0, func(c mqtt.Client, m mqtt.Message) {
		responseChan <- string(m.Payload())
	})

	params := MQTTParams{
		client:   client,
		clientID: clientId,
	}

	mClient.HandleCommand(args, params)
}
