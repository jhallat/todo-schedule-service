package schedule

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func SetupTaskListener(url string) {

	go func() {
		conn, err := amqp.Dial(url)
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

	go func() {
		conn, err := amqp.Dial(url)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		channel, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer channel.Close()

		messages, err := channel.Consume(
			"schedule.item.completed", // queue
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
				var maxReached UpdatedWeeklyMaxReached
				json.Unmarshal(message.Body, &maxReached)
				err := updateWeeklyMaxReached(maxReached.TaskId, maxReached.MaxReached)
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
