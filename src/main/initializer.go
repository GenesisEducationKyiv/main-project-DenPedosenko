package main

import (
	"ses.genesis.com/exchange-web-service/src/notification"
	"ses.genesis.com/exchange-web-service/src/persistent"
	"ses.genesis.com/exchange-web-service/src/service"
)

func initialize() service.InternalService {
	notificationService := notification.NewEmailSender()
	persistentService := persistent.NewFileStorage(persistent.NewFileProcessor())
	externalService := service.NewExternalExchangeAPIController()

	return service.NewMainService(externalService, persistentService, notificationService)
}
