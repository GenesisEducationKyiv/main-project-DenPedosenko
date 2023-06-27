package main

import (
	"github.com/go-resty/resty/v2"
	"ses.genesis.com/exchange-web-service/src/main/config"
	"ses.genesis.com/exchange-web-service/src/main/notification"
	"ses.genesis.com/exchange-web-service/src/main/persistent"
	"ses.genesis.com/exchange-web-service/src/main/service"
	"ses.genesis.com/exchange-web-service/src/main/service/errormapper"
)

const (
	// ConfigPath is a path to config file
	configPath      = "src/main/resources/application.yaml"
	fileStoragePath = "src/main/resources/emails.txt"
)

func initialize() service.InternalService {
	configLoader := config.NewConfigLoader(configPath)

	ctx, err := configLoader.GetContext()

	if err != nil {
		panic(err)
	}

	notificationService := notification.NewEmailSender(ctx, notification.NewSMTPProtocolService())
	persistentService := persistent.NewFileStorage(persistent.NewFileProcessor(fileStoragePath))
	externalService := service.NewExternalExchangeAPIController(config.GetConfigFromContext(ctx), resty.New())
	storageToHTTPMapper := errormapper.NewStorageErrorToHTTPMapper()

	return service.NewMainService(externalService, persistentService, notificationService, storageToHTTPMapper)
}
