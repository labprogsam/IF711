package main

import (
	"atividade-5/base"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaBroker  = "localhost:9092"
	requestTopic = "file-requests"
	clientDir    = "../files_client"
)

type ClientKafka struct {
	base *base.BaseClient
}

func NewClientKafka() *ClientKafka {
	c := &ClientKafka{}
	c.base = base.NewBaseClient(c)
	return c
}

type KafkaParams struct {
	reader   *kafka.Reader
	clientId string
}

func (c *ClientKafka) HandleCommand(cmd []string, params KafkaParams) {
	c.base.HandleCommand(cmd, params)
}

func (c *ClientKafka) SendCommand(command string, input any) string {
	params, ok := input.(KafkaParams)
	if !ok {
		fmt.Println("parametros errados")
	}
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaBroker},
		Topic:   requestTopic,
	})
	defer writer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	writer.WriteMessages(ctx, kafka.Message{
		Value: []byte(command),
		Headers: []kafka.Header{
			{Key: "timestamp", Value: []byte(time.Now().Format("2006-01-02 15:04:05.999999999 -0700 MST"))},
			{Key: "clientID", Value: []byte(params.clientId)},
		},
	})

	// Listen for response
	msg, err := params.reader.ReadMessage(context.Background())
	if err != nil {
		panic(err)
	}

	return string(msg.Value)
}

func main() {
	kClient := NewClientKafka()

	// Retrieve Arguments
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Usage: <LIST | UPLOAD | DOWNLOAD> <ID> <FILENAME>")
		return
	}
	clientID := args[1]

	// Kafka Reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   "file-responses-" + clientID,
		GroupID: clientID + "-group",
	})
	defer reader.Close()

	params := KafkaParams{
		reader:   reader,
		clientId: clientID,
	}

	kClient.HandleCommand(args, params)
}
