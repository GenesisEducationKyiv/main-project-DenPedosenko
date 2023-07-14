package exchange

import (
	"container/list"
	"errors"
	"exchange-web-service/main/domain/config"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

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

func (s *Service) CurrentRate(from, to string) (float64, error) {
	return s.getRate(s.externalAPIs.Front(), from, to)
}

func (s *Service) getRate(val *list.Element, from, to string) (float64, error) {
	if val == nil {
		logrus.Error("No exchange_service API available")
		return 0, errors.New("no exchange_service API available")
	}

	api, ok := val.Value.(RateAPI)

	if !ok {
		logrus.Error("Can't get rateApi from chain")
		return 0, errors.New("can't get rateApi from chain")
	}

	rate, err := api.GetRate(from, to)
	if err != nil {
		return s.getRate(val.Next(), from, to)
	}

	return rate, nil
}
