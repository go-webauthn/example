package main

import (
	"github.com/go-webauthn/example/internal/configuration"
	"github.com/go-webauthn/example/internal/logging"
	"github.com/go-webauthn/example/internal/server"
)

func main() {
	config, err := configuration.Load([]string{"config.yaml"}, true, nil)
	if err != nil {
		panic(err)
	}

	logger, err := logging.Configure(&config.Log)
	if err != nil {
		panic(err)
	}

	logger.Info("server is starting")

	if err = server.Run(config); err != nil {
		panic(err)
	}
}
