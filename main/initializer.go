package main

import (
	"github.com/go-resty/resty/v2"
	config2 "ses.genesis.com/exchange-web-service/main/config"
	notification2 "ses.genesis.com/exchange-web-service/main/notification"
	persistent2 "ses.genesis.com/exchange-web-service/main/persistent"
	"ses.genesis.com/exchange-web-service/main/service"
	"ses.genesis.com/exchange-web-service/main/service/errormapper"
)

const (
	configPath      = "main/resources/application.yaml"
	fileStoragePath = "main/resources/emails.txt"
)

func initialize() service.InternalService {
	configLoader := config2.NewConfigLoader(configPath)

	ctx, err := configLoader.GetContext()

	if err != nil {
		panic(err)
	}

	notificationService := notification2.NewEmailSender(ctx, notification2.NewSMTPProtocolService())
	persistentService := persistent2.NewFileStorage(persistent2.NewFileProcessor(fileStoragePath))
	externalService := service.NewExternalExchangeAPIController(config2.GetConfigFromContext(ctx), resty.New())
	storageToHTTPMapper := errormapper.NewStorageErrorToHTTPMapper()

	return service.NewMainService(externalService, persistentService, notificationService, storageToHTTPMapper)
}
