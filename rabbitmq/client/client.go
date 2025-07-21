package main

import (
	"atividade-5/base"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

const clientFiles = "../files_client"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type ClientRabbitMQ struct {
	base *base.BaseClient
}

func NewClientRabbitMQ() *ClientRabbitMQ {
	c := &ClientRabbitMQ{}
	c.base = base.NewBaseClient(c)
	return c
}

type RabbitParams struct {
	ch         *amqp.Channel
	replies    <-chan amqp.Delivery
	corrId     string
	replyQueue string
}

func (c *ClientRabbitMQ) SendCommand(command string, input any) string {
	params, ok := input.(RabbitParams)
	if !ok {
		fmt.Println("parametros errados")
	}

	err := params.ch.Publish(
		"", "file/request", false, false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: params.corrId,
			ReplyTo:       params.replyQueue,
			Body:          []byte(command),
		})
	failOnError(err, "Failed to publish message")

	for d := range params.replies {
		if d.CorrelationId == params.corrId {
			return string(d.Body)
		}
	}
	return ""
}

func (c *ClientRabbitMQ) HandleCommand(cmd []string, params RabbitParams) {
	c.base.HandleCommand(cmd, params)
}

func main() {
	os.MkdirAll(clientFiles, os.ModePerm)
	rClient := NewClientRabbitMQ()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	failOnError(err, "Failed to declare a reply queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Usage: <LIST | UPLOAD | DOWNLOAD> <ID> <FILENAME>")
		return
	}

	corrId := args[1]

	params := RabbitParams{
		ch:         ch,
		replies:    msgs,
		corrId:     corrId,
		replyQueue: q.Name,
	}

	rClient.HandleCommand(args, params)
}
