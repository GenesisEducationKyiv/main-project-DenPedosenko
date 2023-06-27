package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"ses.genesis.com/exchange-web-service/main/config"

	"github.com/go-resty/resty/v2"
)

type ExchangeRateResponse struct {
	Rate float64 `json:"rate"`
}

type ExternalExchangeAPIController struct {
	config *config.AppConfig
	client *resty.Client
}

func NewExternalExchangeAPIController(conf *config.AppConfig, client *resty.Client) *ExternalExchangeAPIController {
	return &ExternalExchangeAPIController{
		config: conf,
		client: client,
	}
}

func (controller *ExternalExchangeAPIController) CurrentBTCToUAHRate() (float64, error) {
	var response ExchangeRateResponse

	resp, err := controller.client.R().
		SetHeader("X-CoinAPI-Key", controller.config.APIKey).
		Get(controller.config.APIURL)

	if err != nil {
		return 0, fmt.Errorf("failed to perform API request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("unexpected API response: %s", resp.Status())
	}

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return 0, fmt.Errorf("failed to parse API response: %w", err)
	}

	return response.Rate, nil
}
