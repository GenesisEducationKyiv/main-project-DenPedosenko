package provider

import (
	"exchange-web-service/domain/config"
	"exchange-web-service/domain/logger"

	"github.com/go-resty/resty/v2"
)

type factory struct {
	client                *resty.Client
	clientLoggerDecorator *logger.Decorator
}

func NewAPIFactory(client *resty.Client, clientLoggerDecorator *logger.Decorator) *factory {
	return &factory{
		client:                client,
		clientLoggerDecorator: clientLoggerDecorator,
	}
}

func (factory *factory) CoinAPIProvider(api config.ConfigAPI) *CoinAPIProvider {
	return NewCoinAPIProvider(&api, factory.clientLoggerDecorator.NewLogResponseDecorator(factory.client))
}

func (factory *factory) CoinGeckoAPIProvider(api config.ConfigAPI) *CoinGeckoProvider {
	return NewCoinGeckoProvider(&api, factory.clientLoggerDecorator.NewLogResponseDecorator(factory.client))
}
func (factory *factory) KuCoinAPIProvider(api config.ConfigAPI) *KuCoinAPIProvider {
	return NewKuCoinProvider(&api, factory.clientLoggerDecorator.NewLogResponseDecorator(factory.client))
}
