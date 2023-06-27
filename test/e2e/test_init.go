package e2e

import (
	"github.com/go-resty/resty/v2"
	"ses.genesis.com/exchange-web-service/main/config"
	"ses.genesis.com/exchange-web-service/main/notification"
	"ses.genesis.com/exchange-web-service/main/persistent"
	"ses.genesis.com/exchange-web-service/main/service"
	"ses.genesis.com/exchange-web-service/main/service/errormapper"
)

const (
	// ConfigPath is a path to config file
	configPath      = "../application.yaml"
	fileStoragePath = "emails.txt"
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
