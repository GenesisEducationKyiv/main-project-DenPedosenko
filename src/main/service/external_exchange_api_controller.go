package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"ses.genesis.com/exchange-web-service/src/main/config"

	"github.com/go-resty/resty/v2"
)

type ExchangeRateResponse struct {
	Rate float64 `json:"rate"`
}

type ExternalExchangeAPIController struct {
	ctx context.Context
}

func NewExternalExchangeAPIController(ctx context.Context) *ExternalExchangeAPIController {
	return &ExternalExchangeAPIController{
		ctx: ctx,
	}
}

func (controller *ExternalExchangeAPIController) GetCurrentBTCToUAHRate() (float64, error) {
	var response ExchangeRateResponse

	conf := config.GetConfig(controller.ctx)
	client := resty.New()
	resp, err := client.R().
		SetHeader("X-CoinAPI-Key", conf.APIKey).
		Get(conf.APIURL)

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
