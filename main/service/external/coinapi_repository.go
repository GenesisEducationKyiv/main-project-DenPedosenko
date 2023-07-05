package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"ses.genesis.com/exchange-web-service/main/config"
)

type RateAPI interface {
	GetRate(from, to string) (float64, error)
}

type CoinAPIRepository struct {
	config *config.ConfigAPI
	client *resty.Client
}

func NewCoinAPIRepository(conf *config.ConfigAPI, client *resty.Client) *CoinAPIRepository {
	return &CoinAPIRepository{
		config: conf,
		client: client,
	}
}

type CoinAPIResponse struct {
	Rate float64 `json:"rate"`
}

func (repository CoinAPIRepository) GetRate(from, to string) (float64, error) {
	var response CoinAPIResponse

	resp, err := repository.client.R().
		SetHeader("X-CoinAPI-Key", repository.config.Key).
		Get(fmt.Sprintf("%s/%s/%s", repository.config.URL, strings.ToUpper(from), strings.ToUpper(to)))

	if err != nil {
		return 0, fmt.Errorf("failed to perform API request: %w", err)
	}

	logrus.Infof("CoinAPI response: %s", resp.String())

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("unexpected API response: %s", resp.Status())
	}

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return 0, fmt.Errorf("failed to parse API response: %w", err)
	}

	return response.Rate, nil
}
