package task

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func SetupListener() {

	go func() {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		channel, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer channel.Close()

		messages, err := channel.Consume(
			"task.description", // queue
			"",                 // consumer
			true,               // auto-ack
			false,              // exclusive
			false,              // no-local
			false,              // no-wait
			nil,                // args
		)
		failOnError(err, "Failed to register a consumer")

		for range time.Tick(time.Second) {
			for message := range messages {
				var updatedTask UpdatedTask
				json.Unmarshal(message.Body, &updatedTask)
				err := updateDescription(updatedTask)
				if err != nil {
					log.Print(err)
				}
			}
		}
	}()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
