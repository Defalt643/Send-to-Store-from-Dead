package main

import amqp "github.com/rabbitmq/amqp091-go"

// RabbitMQ variables
var (
	WatchChan        interface{}
	RabbitCloseError chan *amqp.Error
)

// Elasticsearch operation
var (
	InsertDocument = "I"
	UpdateDocuemnt = "U"
)

// Graceful shutdown variable
var (
	IsNodeShutdown = false
)
