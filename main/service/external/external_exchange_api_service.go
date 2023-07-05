package external

import (
	"container/list"
	"errors"

	"github.com/sirupsen/logrus"
	"ses.genesis.com/exchange-web-service/main/config"

	"github.com/go-resty/resty/v2"
)

type RateService interface {
	CurrentRate(from string, to string) (float64, error)
}

type RateAPI interface {
	GetRate(from, to string) (float64, error)
}

type Service struct {
	config       *config.AppConfig
	client       *resty.Client
	externalAPIs *list.List
}

func NewExternalExchangeAPIService(conf *config.AppConfig, client *resty.Client, apis *list.List) *Service {
	return &Service{
		config:       conf,
		client:       client,
		externalAPIs: apis,
	}
}

func (controller *Service) CurrentRate(from, to string) (float64, error) {
	return getRate(controller.externalAPIs.Front(), from, to)
}

func getRate(val *list.Element, from, to string) (float64, error) {
	if val == nil {
		logrus.Error("No external API available")
		return 0, errors.New("no external API available")
	}

	api, ok := val.Value.(RateAPI)

	if !ok {
		logrus.Error("Can't get rateApi from chain")
		return 0, errors.New("can't get rateApi from chain")
	}

	rate, err := api.GetRate(from, to)
	if err != nil {
		return getRate(val.Next(), from, to)
	}

	return rate, nil
}
