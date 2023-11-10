package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"smq/pkg/services"
	"smq/pkg/smqerrors"
	"syscall"
)

func main() {
	service, err := services.NewSMQService(services.WithViperCfg(), services.WithSimpleMessageQueue())
	smqerrors.Panic(err)

	done := make(chan bool, 1)

	signalsChannel := make(chan os.Signal, 1)
	signal.Notify(signalsChannel, syscall.SIGINT, syscall.SIGTERM)

	_, cancel := context.WithCancel(context.Background())

	go func() {
		service.Start()
	}()

	go func() {
		log.Println("waiting for termination signal")
		sig := <-signalsChannel
		log.Println("got termination signal", sig)
		//http.Post(ntfyURL, "text/plain",
		//	strings.NewReader("Service Old autovit STOPPED 😡"))
		cancel()
		done <- true
	}()

	<-done
}
