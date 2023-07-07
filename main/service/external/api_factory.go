package external

import (
	"github.com/go-resty/resty/v2"
	"ses.genesis.com/exchange-web-service/main/config"
	"ses.genesis.com/exchange-web-service/main/logger"
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

func (factory *factory) CoinAPIRepository(api config.ConfigAPI) *CoinAPIProvider {
	return NewCoinAPIProvider(&api, factory.clientLoggerDecorator.NewLogResponseDecorator(factory.client))
}

func (factory *factory) CoinGeckoAPIRepository(api config.ConfigAPI) *CoinGeckoProvider {
	return NewCoinGeckoProvider(&api, factory.clientLoggerDecorator.NewLogResponseDecorator(factory.client))
}
func (factory *factory) KuCoinAPIRepository(api config.ConfigAPI) *KuCoinAPIProvider {
	return NewKuCoinProvider(&api, factory.clientLoggerDecorator.NewLogResponseDecorator(factory.client))
}
