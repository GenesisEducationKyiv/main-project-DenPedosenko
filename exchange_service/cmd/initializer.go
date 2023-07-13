package cmd

import (
	"container/list"

	"exchange-web-service/service"

	"exchange-web-service/application/exchange"
	"exchange-web-service/application/exchange/provider"
	"exchange-web-service/application/notification"
	"exchange-web-service/persistent"

	"exchange-web-service/domain/config"
	"exchange-web-service/domain/logger"
	"exchange-web-service/presentation/handler"
	"exchange-web-service/presentation/handler/errormapper"

	"github.com/go-resty/resty/v2"
)

const (
	configPath      = "resources/application.yaml"
	fileStoragePath = "resources/emails.txt"
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
	configLoader := config.NewConfigLoader(configPath)

	ctx, err := configLoader.GetContext()
	if err != nil {
		panic(err)
	}

	conf := config.GetConfigFromContext(ctx)
	notificationService := notification.NewEmailSender(ctx, notification.NewSMTPProtocolService())
	persistentService := persistent.NewFileStorage(persistent.NewFileProcessor(fileStoragePath))
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
