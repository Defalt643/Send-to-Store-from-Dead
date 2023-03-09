package main

import (
	"flag"

	"veteran.socialenable.co/se4/middlelibrary/graylog2"
)

// Env is global variable for set environment
var Env string

// BulkSize is global variable for elasticsearch bulk size
var BulkSize int

// LogLevel default is debug
var LogLevel string

// EnableGraylog default is false
var EnableGraylog bool

// QueueConsumeName is queue name
var QueueConsumeName string

// QueuePublishProcessName is queue name
var QueuePublishDeadName string

func init() {
	flag.StringVar(&Env, "env", "beta", "Set environment")
	flag.StringVar(&LogLevel, "log-level", "debug", "Set log level")
	flag.BoolVar(&EnableGraylog, "enable-graylog", false, "Set graylog enabling")
	flag.IntVar(&BulkSize, "bulk-size", 1, "Set bulk size")
	flag.StringVar(&QueueConsumeName, "queue-consume-name", "ha_oplog.speech.recognition.youtube", "Queue for consume data")
	flag.StringVar(&QueuePublishDeadName, "queue-publish-dead-name", "ha_oplog.speech.recognition.youtube.dead", "Queue for publish not found data")
	flag.Parse()
}

func main() {

	// Get configuration
	configuration := GetConfiguration()

	graylog2config := graylog2.Graylog2{
		Version:     "v1",
		Address:     configuration.Graylog.Address,
		Enable:      EnableGraylog,
		LevelString: LogLevel,
	}

	err := graylog2.New(graylog2config)
	if err != nil {
		graylog2.SetEnable(false)
		graylog2.Errorf(13400, "Connect graylog error : %s", err)
	}

	// Rabbitmq initializer
	configuration.Rabbit.Connection = configuration.Rabbit.Connect()
	defer configuration.Rabbit.Connection.Close()
	channel, err := configuration.Rabbit.Connection.Channel()
	if err != nil {
		graylog2.Panicf(13400, "Cannot create channel of rabbitmq: %s", err)
	}

	configuration.Rabbit.Channel = channel
	defer configuration.Rabbit.Channel.Close()

	configuration.Rabbit.ReconnectRabbitMQ()

	resource := Resource{
		Rabbit: configuration.Rabbit,
	}

	go resource.WatchShutdownSignal()

	resource.Watch()
}
