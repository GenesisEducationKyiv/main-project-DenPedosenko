package e2e

import (
	"container/list"
	"exchange-web-service/application/exchange"
	"exchange-web-service/application/exchange/provider"
	"exchange-web-service/application/notification"
	"exchange-web-service/domain/config"
	"exchange-web-service/domain/logger"
	"exchange-web-service/persistent"
	"exchange-web-service/presentation/handler"
	"exchange-web-service/presentation/handler/errormapper"
	"exchange-web-service/service"

	"github.com/go-resty/resty/v2"
)

const (
	configPath      = "../application.yaml"
	fileStoragePath = "emails.txt"
)

func initialize() *service.MainService {
	configLoader := config.NewConfigLoader(configPath)

	ctx, err := configLoader.GetContext()
	testLogger := TestLogger{}

	if err != nil {
		panic(err)
	}

	notificationService := notification.NewEmailSender(ctx, notification.NewSMTPProtocolService(), testLogger)
	persistentService := persistent.NewFileStorage(persistent.NewFileProcessor(fileStoragePath, testLogger), testLogger)
	apis := list.New()
	conf := config.GetConfigFromContext(ctx)
	apisFactory := provider.NewAPIFactory(resty.New(), logger.NewLogger())

	apis.PushBack(apisFactory.CoinAPIProvider(conf.CoinAPI))
	apis.PushBack(apisFactory.CoinGeckoAPIProvider(conf.CoinGecko))
	apis.PushBack(apisFactory.KuCoinAPIProvider(conf.KuCoin))

	externalService := exchange.NewExternalExchangeAPIService(config.GetConfigFromContext(ctx), resty.New(), apis, testLogger)

	storageToHTTPMapper := errormapper.NewStorageErrorToHTTPMapper()

	rateController := handler.NewRateHandler(externalService, testLogger)
	emailController := handler.NewEmailHandler(persistentService, storageToHTTPMapper, testLogger)
	notificationController := handler.NewNotificationHandler(externalService, notificationService, persistentService, testLogger)

	return service.NewMainService(rateController, emailController, notificationController)
}
