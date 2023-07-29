package provider

import (
	"encoding/json"
	"exchange-web-service/domain/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

type CoinAPIProvider struct {
	config *config.ConfigAPI
	client *resty.Client
}

func NewCoinAPIProvider(conf *config.ConfigAPI, client *resty.Client) *CoinAPIProvider {
	return &CoinAPIProvider{
		config: conf,
		client: client,
	}
}

type CoinAPIResponse struct {
	Rate float64 `json:"rate"`
}

func (repository CoinAPIProvider) GetRate(from, to string) (float64, error) {
	var response CoinAPIResponse

	resp, err := repository.client.R().
		SetHeader("X-CoinAPI-Key", repository.config.Key).
		Get(fmt.Sprintf("%s/%s/%s", repository.config.URL, strings.ToUpper(from), strings.ToUpper(to)))

	if err != nil {
		return 0, fmt.Errorf("failed to perform API request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("unexpected API response: %s", resp.Status())
	}

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return 0, fmt.Errorf("failed to parse API Data: %w", err)
	}

	return response.Rate, nil
}
