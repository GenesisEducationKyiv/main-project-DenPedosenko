package external_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"ses.genesis.com/exchange-web-service/main/config"
	"ses.genesis.com/exchange-web-service/main/service/external"
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
				response := external.CoinAPIResponse{Rate: 500000}
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
			conf := &config.ConfigAPI{
				URL: scenario.server.URL,
				Key: "test-key",
			}
			repository := external.NewCoinAPIRepository(conf, client)

			rate, err := repository.GetRate("btc", "uah")

			assert.Equal(t, scenario.expErr, err)
			assert.Equal(t, scenario.expected, rate)
		})
	}
}
