package cmd

import (
	"container/list"
	"context"
	"exchange-web-service/application/exchange"
	"exchange-web-service/application/exchange/provider"
	"exchange-web-service/application/notification"
	"exchange-web-service/domain/config"
	"exchange-web-service/domain/logger"
	"exchange-web-service/domain/rabittmq"
	"exchange-web-service/persistent"
	"exchange-web-service/presentation/handler"
	"exchange-web-service/presentation/handler/errormapper"
	"exchange-web-service/service"

	"github.com/go-resty/resty/v2"
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

func (i MainInitializer) initialize(ctx context.Context, rabbitMQLogger *rabittmq.Logger) (*service.MainService, *service.SagaService) {
	conf := config.GetConfigFromContext(ctx)
	notificationService := notification.NewEmailSender(ctx, notification.NewSMTPProtocolService(), rabbitMQLogger)
	persistentService := persistent.NewFileStorage(persistent.NewFileProcessor(fileStoragePath, rabbitMQLogger), rabbitMQLogger)
	apisFactory := provider.NewAPIFactory(resty.New(), logger.NewLogger())

	apis := list.New()
	apis.PushFront(apisFactory.CoinGeckoAPIProvider(conf.CoinGecko))
	apis.PushFront(apisFactory.CoinAPIProvider(conf.CoinAPI))
	apis.PushFront(apisFactory.KuCoinAPIProvider(conf.KuCoin))

	externalService := exchange.NewExternalExchangeAPIService(config.GetConfigFromContext(ctx), resty.New(), apis, rabbitMQLogger)
	storageToHTTPMapper := errormapper.NewStorageErrorToHTTPMapper()

	rateHandler := handler.NewRateHandler(externalService, rabbitMQLogger)
	emailHandler := handler.NewEmailHandler(persistentService, storageToHTTPMapper, rabbitMQLogger)
	notificationHandler := handler.NewNotificationHandler(externalService, notificationService, persistentService, rabbitMQLogger)

	return service.NewMainService(rateHandler, emailHandler, notificationHandler),
		service.NewSagaService(persistentService, storageToHTTPMapper)
}
