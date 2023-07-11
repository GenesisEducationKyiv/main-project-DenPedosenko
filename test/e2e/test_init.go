package e2e

import (
	"container/list"

	"ses.genesis.com/exchange-web-service/main/application/exchange"
	"ses.genesis.com/exchange-web-service/main/application/exchange/provider"
	"ses.genesis.com/exchange-web-service/main/application/notification"
	"ses.genesis.com/exchange-web-service/main/persistent"

	"github.com/go-resty/resty/v2"
	"ses.genesis.com/exchange-web-service/main/domain/config"
	"ses.genesis.com/exchange-web-service/main/domain/logger"
	"ses.genesis.com/exchange-web-service/main/presentation/handler"
	"ses.genesis.com/exchange-web-service/main/presentation/handler/errormapper"
	"ses.genesis.com/exchange-web-service/main/service"
)

const (
	configPath      = "../application.yaml"
	fileStoragePath = "emails.txt"
)

func initialize() *service.MainService {
	configLoader := config.NewConfigLoader(configPath)

	ctx, err := configLoader.GetContext()

	if err != nil {
		panic(err)
	}

	notificationService := notification.NewEmailSender(ctx, notification.NewSMTPProtocolService())
	persistentService := persistent.NewFileStorage(persistent.NewFileProcessor(fileStoragePath))
	apis := list.New()
	conf := config.GetConfigFromContext(ctx)
	apisFactory := provider.NewAPIFactory(resty.New(), logger.NewLogger())

	apis.PushBack(apisFactory.CoinAPIProvider(conf.CoinAPI))
	apis.PushBack(apisFactory.CoinGeckoAPIProvider(conf.CoinGecko))
	apis.PushBack(apisFactory.KuCoinAPIProvider(conf.KuCoin))

	externalService := exchange.NewExternalExchangeAPIService(config.GetConfigFromContext(ctx), resty.New(), apis)

	storageToHTTPMapper := errormapper.NewStorageErrorToHTTPMapper()

	rateController := handler.NewRateHandler(externalService)
	emailController := handler.NewEmailHandler(persistentService, storageToHTTPMapper)
	notificationController := handler.NewNotificationHandler(externalService, notificationService, persistentService)

	return service.NewMainService(rateController, emailController, notificationController)
}
