package main

import (
	"atividade-5/base"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

const (
	queueName = "file/request"
	serverDir = "../files_server"
)

type ServerRabbitMQ struct {
	*base.BaseServer
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func handleRequest(d amqp.Delivery, ch *amqp.Channel) {
	payload := string(d.Body)
	server := &ServerRabbitMQ{&base.BaseServer{}}

	parts := strings.SplitN(payload, " ", 3)

	if len(parts) < 1 {
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

	response := server.HandleRequest(command, filename, data)

	err := ch.Publish(
		"",        // exchange
		d.ReplyTo, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: d.CorrelationId,
			Body:          []byte(response),
		})
	failOnError(err, "Failed to publish a response")
}

func main() {
	os.MkdirAll(serverDir, os.ModePerm)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			handleRequest(d, ch)
		}
	}()

	fmt.Println("AMQP server is running")
	<-forever
}
