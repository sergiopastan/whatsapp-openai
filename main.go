package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/sergiopastan/whatsapp-openai/config"
	"github.com/sergiopastan/whatsapp-openai/database"
	"github.com/sergiopastan/whatsapp-openai/openai"
	"github.com/sergiopastan/whatsapp-openai/whatsapp"
)

func main() {
	var conf = config.Load()

	configureLogger()
	db, err := database.Connect(conf.DbConfig)
	if err != nil {
		log.WithError(err).Fatal("failed to create db connection")
	}
	wspClient, err := whatsapp.NewClient(db)
	openaiHandler := openai.NewHandler(wspClient)
	wspClient.AddEventHandler(whatsapp.MessageReceiptHandler(openaiHandler))
	if err != nil {
		log.WithError(err).Fatal("failed to create whatsapp client")
	}
	err = wspClient.Start(context.Background())
	if err != nil {
		log.WithError(err).Fatal("failed to start whatsapp connection")
	}

	waitForInterrupt()

	log.Info("cleaning up...")

	wspClient.Disconnect()

	log.Info("successful shutdown")
}

func configureLogger() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{DisableColors: true, FullTimestamp: true})
}

func waitForInterrupt() {
	osSignals := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-osSignals // This blocks the main thread until an interrupt or termination is received
}
