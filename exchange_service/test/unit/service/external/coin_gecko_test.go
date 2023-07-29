package external_test

import (
	"errors"
	"exchange-web-service/application/exchange/provider"
	"exchange-web-service/domain/config"
	"exchange-web-service/domain/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestGetRateFromGecko(t *testing.T) {
	for _, scenario := range []Scenario{
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
			client := resty.New().SetBaseURL(scenario.server.URL)
			client = logger.NewLogger().NewLogResponseDecorator(client)

			conf := &config.ConfigAPI{
				URL: scenario.server.URL,
			}
			repository := provider.NewCoinGeckoProvider(conf, client)
			test(t, scenario, repository)
		})
	}
}
