package cmd

import (
	"container/list"

	"ses.genesis.com/exchange-web-service/main/service"

	"ses.genesis.com/exchange-web-service/main/application/exchange"
	"ses.genesis.com/exchange-web-service/main/application/exchange/provider"
	"ses.genesis.com/exchange-web-service/main/application/notification"
	persistent2 "ses.genesis.com/exchange-web-service/main/persistent"

	"github.com/go-resty/resty/v2"
	"ses.genesis.com/exchange-web-service/main/domain/config"
	"ses.genesis.com/exchange-web-service/main/domain/logger"
	"ses.genesis.com/exchange-web-service/main/presentation/handler"
	"ses.genesis.com/exchange-web-service/main/presentation/handler/errormapper"
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
	persistentService := persistent2.NewFileStorage(persistent2.NewFileProcessor(fileStoragePath))
	apisFactory := provider.NewAPIFactory(resty.New(), logger.NewLogger())

	apis := list.New()
	apis.PushFront(apisFactory.CoinGeckoAPIProvider(conf.CoinGecko))
	apis.PushFront(apisFactory.CoinAPIProvider(conf.CoinAPI))
	apis.PushFront(apisFactory.KuCoinAPIProvider(conf.KuCoin))

	externalService := exchange.NewExternalExchangeAPIService(config.GetConfigFromContext(ctx), resty.New(), apis)
	storageToHTTPMapper := errormapper.NewStorageErrorToHTTPMapper()

	rateHandler := handler.NewRateHandler(externalService)
	emailHandler := handler.NewEmailHandler(persistentService, storageToHTTPMapper)
	notificationHandler := handler.NewNotificationHandler(externalService, notificationService, persistentService)

	return service.NewMainService(rateHandler, emailHandler, notificationHandler)
}
