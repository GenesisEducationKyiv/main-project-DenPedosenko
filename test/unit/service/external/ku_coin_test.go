package external_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"ses.genesis.com/exchange-web-service/main/logger"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"ses.genesis.com/exchange-web-service/main/config"
	"ses.genesis.com/exchange-web-service/main/service/external"
)

type Scenario struct {
	name     string
	server   *httptest.Server
	expected float64
	expErr   error
}

func TestGetRateFromKuCoinApi(t *testing.T) {
	for _, scenario := range []Scenario{
		{
			name: "success",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := external.KuCoinResponse{
					Code: "200",
					Data: map[string]string{"btc": "500000"},
				}
				_ = json.NewEncoder(w).Encode(response)
			})),
			expected: 500000,
			expErr:   nil,
		},
		{
			name: "server error",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "internal server error", http.StatusInternalServerError)
			})),
			expected: 0,
			expErr:   errors.New("unexpected API response: 500 Internal Server Error"),
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			client := resty.New().SetBaseURL(scenario.server.URL)
			client = logger.NewLogger().NewLogResponseDecorator(client)

			conf := &config.ConfigAPI{
				URL: scenario.server.URL,
			}
			repository := external.NewKuCoinProvider(conf, client)
			test(t, scenario, repository)
		})
	}
}

func test(t *testing.T, scenario Scenario, repository external.RateAPI) {
	defer scenario.server.Close()

	rate, err := repository.GetRate("btc", "uah")

	assert.Equal(t, scenario.expErr, err)
	assert.Equal(t, scenario.expected, rate)
}
