package cmd

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
	configPath      = "main/resources/application.yaml"
	fileStoragePath = "emails.txt"
)

type Router[T any] interface {
	CreateRoutes() T
}

type Initializer interface {
	initialize() *service.MainService
}

type MainInitializer struct {
}

func NewInitializer() *MainInitializer {
	return &MainInitializer{}
}

func (i MainInitializer) initialize() *service.MainService {
	configLoader := config2.NewConfigLoader(configPath)

	ctx, err := configLoader.GetContext()
	if err != nil {
		panic(err)
	}

	conf := config2.GetConfigFromContext(ctx)
	notificationService := notification2.NewEmailSender(ctx, notification2.NewSMTPProtocolService())
	persistentService := persistent2.NewFileStorage(persistent2.NewFileProcessor(fileStoragePath))
	apisFactory := provider.NewAPIFactory(resty.New(), logger.NewLogger())

	apis := list.New()
	apis.PushFront(apisFactory.CoinGeckoAPIProvider(conf.CoinGecko))
	apis.PushFront(apisFactory.CoinAPIProvider(conf.CoinAPI))
	apis.PushFront(apisFactory.KuCoinAPIProvider(conf.KuCoin))

	externalService := exchange.NewExternalExchangeAPIService(config2.GetConfigFromContext(ctx), resty.New(), apis)
	storageToHTTPMapper := errormapper.NewStorageErrorToHTTPMapper()

	rateHandler := handler2.NewRateHandler(externalService)
	emailHandler := handler2.NewEmailHandler(persistentService, storageToHTTPMapper)
	notificationHandler := handler2.NewNotificationHandler(externalService, notificationService, persistentService)

	return service.NewMainService(rateHandler, emailHandler, notificationHandler)
}
