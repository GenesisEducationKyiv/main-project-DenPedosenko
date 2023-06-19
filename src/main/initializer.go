package main

import (
	"ses.genesis.com/exchange-web-service/src/config"
	"ses.genesis.com/exchange-web-service/src/notification"
	"ses.genesis.com/exchange-web-service/src/persistent"
	"ses.genesis.com/exchange-web-service/src/service"
)

const (
	// ConfigPath is a path to config file
	configPath = "src/config/application.yaml"
)

func initialize() service.InternalService {
	configLoader := config.NewConfigLoader(configPath)

	ctx, err := configLoader.GetContext()

	if err != nil {
		panic(err)
	}

	notificationService := notification.NewEmailSender(ctx)
	persistentService := persistent.NewFileStorage(persistent.NewFileProcessor())
	externalService := service.NewExternalExchangeAPIController(ctx)

	return service.NewMainService(externalService, persistentService, notificationService)
}
