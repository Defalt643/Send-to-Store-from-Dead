package main

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"veteran.socialenable.co/se4/middlelibrary/graylog2"
)

func (resource Resource) ProcessData(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		if IsNodeShutdown {
			resource.ShutdownSignalProcess(false)
		}
		var message IncommingMessage
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			graylog2.Error(00000, "Unable to unmarshal because: ", err)
			msg.Nack(false, true)
		}
		resource.Rabbit.publish(message, QueueConsumeName)
		graylog2.Info("[*] message has been sent to ha_oplog.speech.recognition.youtube", message)
		msg.Ack(false)
	}
}
