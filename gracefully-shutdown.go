package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"veteran.socialenable.co/se4/middlelibrary/graylog2"
)

// GracefulShutdown waits for termination syscalls and doing clean up operations after received it
func (resource Resource) WatchShutdownSignal() {
	graylog2.Errorf(9999, "watching for termination signals")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	// when we get a signal, flip the global ShuttingDown flag
	sig := <-sigChan
	graylog2.Errorf(9999, "got signal:", sig)

	// wait for the liveness checks to fail and kubernetes to reconfigure
	graylog2.Errorf(9999, "graceful shutdown has begun")
	IsNodeShutdown = true

	resource.ShutdownSignalProcess(true)
}

func (resource Resource) ShutdownSignalProcess(isSleep bool) {
	graylog2.Error(9999, "Ready to shutdown ")
	if isSleep {
		time.Sleep(time.Second * 30)
	}

	resource.Rabbit.Channel.Close()
	resource.Rabbit.Connection.Close()
	graylog2.Infof("shutdown signal success ...")
	os.Exit(0)
}

func ShutdownSignalProcessInRabbitTask(rabbit Rabbit) {
	graylog2.Error(13400, "Ready to shutdown ShutdownSignalProcessInRabbitTask ")

	rabbit.Channel.Close()
	rabbit.Connection.Close()

	graylog2.Error(13400, "shutdown signal ShutdownSignalProcessInRabbitTask success ...")
	os.Exit(0)
}
