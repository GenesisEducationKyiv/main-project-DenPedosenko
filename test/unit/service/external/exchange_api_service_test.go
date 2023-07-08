package external_test

import (
	"container/list"
	"errors"
	"testing"

	"ses.genesis.com/exchange-web-service/main/application/exchange"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"ses.genesis.com/exchange-web-service/main/domain/config"
	"ses.genesis.com/exchange-web-service/main/domain/logger"
)

type stub struct {
	rate  float64
	err   error
	calls int
}

type apiExpects struct {
	rate float64
	err  error
}

func (m *stub) GetRate(_, _ string) (float64, error) {
	m.calls++
	return m.rate, m.err
}

func createAPIList(api1, api2, api3 apiExpects) *list.List {
	mockAPI1 := &stub{rate: api1.rate, err: api1.err}
	mockAPI2 := &stub{rate: api2.rate, err: api2.err}
	mockAPI3 := &stub{rate: api3.rate, err: api3.err}

	apis := list.New()
	apis.PushBack(mockAPI1)
	apis.PushBack(mockAPI2)
	apis.PushBack(mockAPI3)

	return apis
}

func TestService_CurrentRate(t *testing.T) {
	type scenario struct {
		name          string
		rate          float64
		err           error
		expectedCalls []int
		apis          *list.List
	}

	for _, scenario := range []scenario{
		{
			name:          "success first api",
			rate:          0.5,
			err:           nil,
			expectedCalls: []int{1, 0, 0},
			apis: createAPIList(
				apiExpects{rate: 0.5, err: nil},
				apiExpects{rate: 0.6, err: nil},
				apiExpects{rate: 0.7, err: nil},
			),
		},
		{
			name:          "success second api",
			rate:          0.6,
			err:           nil,
			expectedCalls: []int{1, 1, 0},
			apis: createAPIList(
				apiExpects{rate: 0.5, err: errors.New("error")},
				apiExpects{rate: 0.6, err: nil},
				apiExpects{rate: 0.7, err: nil},
			),
		},
		{
			name:          "success third api",
			rate:          0.7,
			err:           nil,
			expectedCalls: []int{1, 1, 1},
			apis: createAPIList(
				apiExpects{rate: 0.5, err: errors.New("error")},
				apiExpects{rate: 0.6, err: errors.New("error")},
				apiExpects{rate: 0.7, err: nil},
			),
		},
		{
			name:          "fail",
			rate:          0.0,
			err:           errors.New("no exchange API available"),
			expectedCalls: []int{1, 1, 1},
			apis: createAPIList(
				apiExpects{rate: 0.5, err: errors.New("error")},
				apiExpects{rate: 0.6, err: errors.New("error")},
				apiExpects{rate: 0.7, err: errors.New("error")},
			),
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			client := logger.NewLogger().NewLogResponseDecorator(resty.New())

			service := exchange.NewExternalExchangeAPIService(&config.AppConfig{}, client, scenario.apis)
			rate, err := service.CurrentRate("USD", "EUR")
			assert.Equal(t, scenario.err, err)
			assert.Equal(t, scenario.rate, rate)
			for e, i := scenario.apis.Front(), 0; e != nil; e, i = e.Next(), i+1 {
				assert.Equal(t, scenario.expectedCalls[i], e.Value.(*stub).calls)
			}
		})
	}
}
