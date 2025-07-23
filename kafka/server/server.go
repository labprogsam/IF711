package main

import (
	"atividade-5/base"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaBroker  = "localhost:9092"
	requestTopic = "file-requests"
	serverDir    = "../files_server"
)

type ServerKafka struct {
	*base.BaseServer
}

func handleRequest(srv *ServerKafka, writer *kafka.Writer, clientID string, message string, headers []kafka.Header) {
	parts := strings.SplitN(message, " ", 4)

	if len(parts) < 1 {
		log.Println("Comando inválido.")
		return // comando inválido
	}

	command := parts[0]
	filename := ""
	if len(parts) >= 2 {
		filename = parts[1]
	}

	data := ""
	if len(parts) == 3 {
		data = parts[2]
	}
	response := srv.HandleRequest(command, filename, data)

	// Send response to client-specific topic
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:     []byte(clientID),
		Value:   []byte(response),
		Topic:   "file-responses-" + clientID,
		Headers: headers,
	})
	if err != nil {
		log.Println("Failed to send response:", err)
	}
}

func main() {
	os.MkdirAll(serverDir, os.ModePerm)
	srv := &ServerKafka{&base.BaseServer{}}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   requestTopic,
		GroupID: "file-server-group",
	})
	defer reader.Close()

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaBroker},
	})
	defer writer.Close()

	fmt.Println("Server listening on Kafka topic:", requestTopic)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Failed to read message:", err)
			continue
		}
		clientId := ""
		for _, h := range msg.Headers {
			if h.Key == "clientID" {
				clientId = string(h.Value)
				break
			}
		}
		go handleRequest(srv, writer, clientId, string(msg.Value), msg.Headers)
	}
}
