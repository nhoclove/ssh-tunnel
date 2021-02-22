package main

import (
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"sshtunnel/internal/config"
	"sshtunnel/internal/mgr"
)

func main() {
	appConfig, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	tunnelMgr := mgr.New(appConfig)
	tunnelMgr.Start()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	tunnelMgr.Shutdown()
}
