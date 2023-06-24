package main

import (
	"ses.genesis.com/exchange-web-service/src/main/config"
	notification2 "ses.genesis.com/exchange-web-service/src/main/notification"
	persistent2 "ses.genesis.com/exchange-web-service/src/main/persistent"
	service2 "ses.genesis.com/exchange-web-service/src/main/service"
)

const (
	// ConfigPath is a path to config file
	configPath      = "src/main/resources/application.yaml"
	fileStoragePath = "src/main/resources/emails.txt"
)

func initialize() service2.InternalService {
	configLoader := config.NewConfigLoader(configPath)

	ctx, err := configLoader.GetContext()

	if err != nil {
		panic(err)
	}

	notificationService := notification2.NewEmailSender(ctx, notification2.NewSMTPProtocolService())
	persistentService := persistent2.NewFileStorage(persistent2.NewFileProcessor(fileStoragePath))
	externalService := service2.NewExternalExchangeAPIController(ctx)

	return service2.NewMainService(externalService, persistentService, notificationService)
}
