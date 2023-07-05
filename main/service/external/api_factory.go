package external

import (
	"github.com/go-resty/resty/v2"
	"ses.genesis.com/exchange-web-service/main/config"
)

type factory struct {
	client *resty.Client
}

func NewAPIFactory(client *resty.Client) *factory {
	return &factory{
		client: client,
	}
}

func (factory *factory) CoinAPIRepository(api config.ConfigAPI) *CoinAPIRepository {
	return NewCoinAPIRepository(&api, factory.client)
}

func (factory *factory) CoinGeckoAPIRepository(api config.ConfigAPI) *CoinGeckoProvider {
	return NewCoinGeckoProvider(&api, factory.client)
}
func (factory *factory) KuCoinAPIRepository(api config.ConfigAPI) *KuCoinAPIProvider {
	return NewKuCoinRepository(&api, factory.client)
}
