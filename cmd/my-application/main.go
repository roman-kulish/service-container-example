package main

import (
	"log"
	"os"

	"github.com/roman-kulish/service-container-example/cmd/my-application/app"
	"github.com/roman-kulish/service-container-example/cmd/my-application/app/config"
	"github.com/roman-kulish/service-container-example/cmd/my-application/app/service"
)

func main() {
	// Parse flags, configuration file or initialise configuration from
	// environment variables.
	cfg, err := config.NewFromEnv()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialise container.
	cnt, err := service.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Container service error: %v", err)
	}

	// Pass config and container to a "real" main().
	if err := app.Run(cfg, cnt); err != nil {
		// app.Run() takes over logging form this point and so
		// error message is not logged.
		os.Exit(1)
	}
	os.Exit(0)
}
