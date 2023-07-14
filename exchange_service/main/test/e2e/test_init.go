package e2e

import (
	"container/list"
	"exchange-web-service/main/application/exchange"
	"exchange-web-service/main/application/exchange/provider"
	notification2 "exchange-web-service/main/application/notification"
	config2 "exchange-web-service/main/domain/config"
	"exchange-web-service/main/domain/logger"
	persistent2 "exchange-web-service/main/persistent"
	handler2 "exchange-web-service/main/presentation/handler"
	"exchange-web-service/main/presentation/handler/errormapper"
	"exchange-web-service/main/service"

	"github.com/go-resty/resty/v2"
)

const (
	configPath      = "../application.yaml"
	fileStoragePath = "emails.txt"
)

func initialize() *service.MainService {
	configLoader := config2.NewConfigLoader(configPath)

	ctx, err := configLoader.GetContext()

	if err != nil {
		panic(err)
	}

	notificationService := notification2.NewEmailSender(ctx, notification2.NewSMTPProtocolService())
	persistentService := persistent2.NewFileStorage(persistent2.NewFileProcessor(fileStoragePath))
	apis := list.New()
	conf := config2.GetConfigFromContext(ctx)
	apisFactory := provider.NewAPIFactory(resty.New(), logger.NewLogger())

	apis.PushBack(apisFactory.CoinAPIProvider(conf.CoinAPI))
	apis.PushBack(apisFactory.CoinGeckoAPIProvider(conf.CoinGecko))
	apis.PushBack(apisFactory.KuCoinAPIProvider(conf.KuCoin))

	externalService := exchange.NewExternalExchangeAPIService(config2.GetConfigFromContext(ctx), resty.New(), apis)

	storageToHTTPMapper := errormapper.NewStorageErrorToHTTPMapper()

	rateController := handler2.NewRateHandler(externalService)
	emailController := handler2.NewEmailHandler(persistentService, storageToHTTPMapper)
	notificationController := handler2.NewNotificationHandler(externalService, notificationService, persistentService)

	return service.NewMainService(rateController, emailController, notificationController)
}
