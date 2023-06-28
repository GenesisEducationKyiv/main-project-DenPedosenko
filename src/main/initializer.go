package main

import (
	"github.com/go-resty/resty/v2"
	"ses.genesis.com/exchange-web-service/src/config"
	"ses.genesis.com/exchange-web-service/src/notification"
	"ses.genesis.com/exchange-web-service/src/persistent"
	"ses.genesis.com/exchange-web-service/src/service"
	"ses.genesis.com/exchange-web-service/src/service/errormapper"
)

const (
	// ConfigPath is a path to config file
	configPath = "src/resources/application.yaml"
)

func initialize() service.InternalService {
	configLoader := config.NewConfigLoader(configPath)

	ctx, err := configLoader.GetContext()

	if err != nil {
		panic(err)
	}

	notificationService := notification.NewEmailSender(ctx)
	persistentService := persistent.NewFileStorage(persistent.NewFileProcessor())
	externalService := service.NewExternalExchangeAPIController(config.GetConfigFromContext(ctx), resty.New())
	storageToHTTPMapper := errormapper.NewStorageErrorToHTTPMapper()

	return service.NewMainService(externalService, persistentService, notificationService, storageToHTTPMapper)
}
