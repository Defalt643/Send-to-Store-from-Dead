package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"veteran.socialenable.co/se4/middlelibrary/graylog2"
)

type Rabbit struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		fmt.Printf("%s: %s", msg, err)
// 		// log.Panicf("%s: %s", msg, err)
// 	}
// }

// Connect is create connection of rabbitmq
func (rabbit Rabbit) Connect() *amqp.Connection {
	connection, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbit.Username, rabbit.Password, rabbit.Host, rabbit.Port))
	if err != nil {
		graylog2.Panicf(13400, "Cannot connect rabbitmq: %s", err)
	}
	return connection
}

// RabbitConnector is re-establish the connection to RabbitMQ in case
// the connection has died
func (rabbit Rabbit) RabbitConnector() {
	var rabbitErr *amqp.Error
	for {
		rabbitErr = <-RabbitCloseError
		if rabbitErr != nil {
			RabbitCloseError = make(chan *amqp.Error)
			rabbit.Connection.NotifyClose(RabbitCloseError)

			if WatchChan != nil {
				graylog2.Error(13400, "WatchChan is not empty")
				WatchChan.(chan string) <- fmt.Sprintf("RabbitMQ conection have problem : %s", rabbitErr.Reason)
			}
		}
	}
}

// ReconnectRabbitMQ is reconnection rabbit
func (rabbit Rabbit) ReconnectRabbitMQ() {
	// create the rabbitmq error channel
	RabbitCloseError = make(chan *amqp.Error)

	// run the callback in a separate thread
	go rabbit.RabbitConnector()

	// establish the rabbitmq connection by sending
	// an error and thus calling the error callback
	RabbitCloseError <- amqp.ErrClosed
}

// Watch job from master
func (resource Resource) Watch() {
	WatchChan = make(chan string)
	graylog2.Info(QueueConsumeName)

	err := resource.Rabbit.Channel.Qos(
		BulkSize, // prefetch count
		0,        // prefetch size
		false,    // global
	)
	if err != nil {
		graylog2.Panicf(13400, "Failed to set QoS: %s", err)
	}
	msgs, err := resource.Rabbit.Channel.Consume(
		QueuePublishDeadName, // queue
		"",                   // consumer
		false,                // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		graylog2.Panicf(13400, "Failed to register a consumer: %s", err)
	}
	go resource.ProcessData(msgs)
	graylog2.Info(msgs)
	graylog2.Info(" [*] Waiting for messages. To exit press CTRL+C")
	for {
		select {
		case watchChan := <-WatchChan.(chan string):
			graylog2.Panic(13400, watchChan)
		}
	}
}

func (rabbit Rabbit) publish(body IncommingMessage, queueName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// graylog2.Info(body)

	message, errJson := json.Marshal(body)
	if errJson != nil {
		graylog2.Errorf(13400, "Unable to marshal Json because %s", errJson.Error())
		panic(errJson)
	}

	err := rabbit.Channel.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         message,
		})
	if err != nil {
		graylog2.Panic(13400, err)
		panic(err)
	} else {
		graylog2.Info(" [x] Sent message to '" + queueName + "' successfully!: " + string(message))
	}
}
