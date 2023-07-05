package external_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"ses.genesis.com/exchange-web-service/main/config"
	"ses.genesis.com/exchange-web-service/main/service/external"
)

func TestGetRateFromGecko(t *testing.T) {
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
				expectedQuery := "ids=bitcoin&vs_currencies=uah"
				if r.URL.RawQuery != expectedQuery {
					http.Error(w, "unexpected request", http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"bitcoin":{"uah":500000}}`))
			})),
			expected: 500000,
			expErr:   nil,
		},
		{
			name: "bad request",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"error":"bad request"}`))
			})),
			expected: 0,
			expErr:   errors.New("unexpected API response: 400 Bad Request"),
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			defer scenario.server.Close()

			client := resty.New().SetBaseURL(scenario.server.URL)

			conf := &config.ConfigAPI{
				URL: scenario.server.URL,
			}

			repository := external.NewCoinGeckoRepository(conf, client)
			rate, err := repository.GetRate("btc", "uah")

			assert.Equal(t, scenario.expErr, err)
			assert.Equal(t, scenario.expected, rate)
		})
	}
}
