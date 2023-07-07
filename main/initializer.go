package main

import (
	"container/list"

	"ses.genesis.com/exchange-web-service/main/logger"

	"github.com/go-resty/resty/v2"
	"ses.genesis.com/exchange-web-service/main/config"
	"ses.genesis.com/exchange-web-service/main/notification"
	"ses.genesis.com/exchange-web-service/main/persistent"
	"ses.genesis.com/exchange-web-service/main/service"
	"ses.genesis.com/exchange-web-service/main/service/errormapper"
	"ses.genesis.com/exchange-web-service/main/service/external"
)

const (
	configPath      = "main/resources/application.yaml"
	fileStoragePath = "main/resources/emails.txt"
)

func initialize() *service.MainService {
	configLoader := config.NewConfigLoader(configPath)

	ctx, err := configLoader.GetContext()
	if err != nil {
		panic(err)
	}

	conf := config.GetConfigFromContext(ctx)
	notificationService := notification.NewEmailSender(ctx, notification.NewSMTPProtocolService())
	persistentService := persistent.NewFileStorage(persistent.NewFileProcessor(fileStoragePath))
	apisFactory := external.NewAPIFactory(resty.New(), logger.NewLogger())

	apis := list.New()
	apis.PushFront(apisFactory.CoinGeckoAPIRepository(conf.CoinGecko))
	apis.PushFront(apisFactory.CoinAPIRepository(conf.CoinAPI))
	apis.PushFront(apisFactory.KuCoinAPIRepository(conf.KuCoin))

	externalService := external.NewExternalExchangeAPIService(config.GetConfigFromContext(ctx), resty.New(), apis)
	storageToHTTPMapper := errormapper.NewStorageErrorToHTTPMapper()

	rateController := service.NewRateController(externalService)
	emailController := service.NewEmailController(persistentService, storageToHTTPMapper)
	notificationController := service.NewNotificationController(externalService, notificationService, persistentService)

	return service.NewMainService(rateController, emailController, notificationController)
}
