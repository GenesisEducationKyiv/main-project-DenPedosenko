package provider

import (
	"encoding/json"
	"exchange-web-service/main/domain/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

type CoinGeckoProvider struct {
	config *config.ConfigAPI
	client *resty.Client
}

func NewCoinGeckoProvider(conf *config.ConfigAPI, client *resty.Client) *CoinGeckoProvider {
	return &CoinGeckoProvider{
		config: conf,
		client: client,
	}
}

var supportedRatesKeys = map[string]string{
	"btc": "bitcoin",
	"eth": "ethereum",
	"uah": "uah",
}

func (repository CoinGeckoProvider) GetRate(from, to string) (float64, error) {
	var response map[string]map[string]float64

	from = supportedRatesKeys[strings.ToLower(from)]

	to = supportedRatesKeys[strings.ToLower(to)]

	resp, err := repository.client.R().
		SetQueryParam("ids", from).
		SetQueryParam("vs_currencies", to).
		Get(repository.config.URL)

	if err != nil {
		return 0, fmt.Errorf("failed to perform API request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("unexpected API response: %s", resp.Status())
	}

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return 0, err
	}

	if err != nil {
		return 0, fmt.Errorf("failed to parse API Data: %w", err)
	}

	price := response[from][to]

	return price, nil
}
