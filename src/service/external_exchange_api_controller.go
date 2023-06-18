package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type ExchangeRateResponse struct {
	Rate float64 `json:"rate"`
}

type ExternalExchangeAPIController struct {
}

func NewExternalExchangeAPIController() *ExternalExchangeAPIController {
	return &ExternalExchangeAPIController{}
}

func (controller *ExternalExchangeAPIController) GetCurrentBTCToUAHRate() (float64, error) {
	var response ExchangeRateResponse

	client := resty.New()
	resp, err := client.R().
		SetHeader("X-CoinAPI-Key", "1840BB94-23AA-4434-B89F-BD0D74FEFB32").
		Get("https://rest.coinapi.io/v1/exchangerate/BTC/UAH")

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
