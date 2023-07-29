package external_test

import (
	"encoding/json"
	"errors"
	"exchange-web-service/application/exchange/provider"
	"exchange-web-service/domain/config"
	"exchange-web-service/domain/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetRateFromCoinApi(t *testing.T) {
	type CoinGeckoScenario struct {
		name     string
		server   *httptest.Server
		expected float64
		expErr   error
	}

	for _, scenario := range []CoinGeckoScenario{
		{
			name: "success",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := provider.CoinAPIResponse{Rate: 500000}
				_ = json.NewEncoder(w).Encode(response)
			})),
			expected: 500000,
			expErr:   nil,
		},
		{
			name: "bad request",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "not found", http.StatusNotFound)
			})),
			expected: 0,
			expErr:   errors.New("unexpected API response: 404 Not Found"),
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			defer scenario.server.Close()

			client := resty.New().SetBaseURL(scenario.server.URL)
			client = logger.NewLogger().NewLogResponseDecorator(client)
			conf := &config.ConfigAPI{
				URL: scenario.server.URL,
				Key: "test-key",
			}
			repository := provider.NewCoinAPIProvider(conf, client)

			rate, err := repository.GetRate("btc", "uah")

			assert.Equal(t, scenario.expErr, err)
			assert.Equal(t, scenario.expected, rate)
		})
	}
}
