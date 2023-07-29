package exchange

import (
	"container/list"
	"errors"
	"exchange-web-service/domain/config"

	"github.com/go-resty/resty/v2"
)

type logger interface {
	Info(string string)
	Error(string string)
	Debug(string string)
	Close()
}

type RateAPI interface {
	GetRate(from, to string) (float64, error)
}

type Service struct {
	config       *config.AppConfig
	client       *resty.Client
	externalAPIs *list.List
	logger       logger
}

func NewExternalExchangeAPIService(conf *config.AppConfig, client *resty.Client, apis *list.List, logger logger) *Service {
	return &Service{
		config:       conf,
		client:       client,
		externalAPIs: apis,
		logger:       logger,
	}
}

func (s *Service) CurrentRate(from, to string) (float64, error) {
	return s.getRate(s.externalAPIs.Front(), from, to)
}

func (s *Service) getRate(val *list.Element, from, to string) (float64, error) {
	if val == nil {
		e := errors.New("no exchange_service API available")
		s.logger.Error(e.Error())

		return 0, e
	}

	api, ok := val.Value.(RateAPI)

	if !ok {
		e := errors.New("can't get rateApi from chain")
		s.logger.Error(e.Error())

		return 0, e
	}

	rate, err := api.GetRate(from, to)
	if err != nil {
		return s.getRate(val.Next(), from, to)
	}

	return rate, nil
}
