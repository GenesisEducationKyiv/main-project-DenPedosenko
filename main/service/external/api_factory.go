package external

import (
	"github.com/go-resty/resty/v2"
	"ses.genesis.com/exchange-web-service/main/config"
)

type APIFactory struct {
	config *config.AppConfig
	client *resty.Client
}

func NewAPIFactory(conf *config.AppConfig, client *resty.Client) *APIFactory {
	return &APIFactory{
		config: conf,
		client: client,
	}
}

func (factory *APIFactory) CoinAPIRepository() *CoinAPIRepository {
	return NewCoinAPIRepository(&factory.config.CoinAPI, factory.client)
}

func (factory *APIFactory) CoinGeckoAPIRepository() *CoinGeckoRepository {
	return NewCoinGeckoRepository(&factory.config.CoinGecko, factory.client)
}
func (factory *APIFactory) KuCoinAPIRepository() *KuCoinAPIRepository {
	return NewKuCoinRepository(&factory.config.KuCoin, factory.client)
}
