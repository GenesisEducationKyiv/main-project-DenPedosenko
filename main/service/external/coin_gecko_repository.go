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

type CoinGeckoRepository struct {
	config *config.ConfigAPI
	client *resty.Client
}

func NewCoinGeckoRepository(conf *config.ConfigAPI, client *resty.Client) *CoinGeckoRepository {
	return &CoinGeckoRepository{
		config: conf,
		client: client,
	}
}

var supportedRatesKeys = map[string]string{
	"btc": "bitcoin",
	"eth": "ethereum",
	"uah": "uah",
}

type coinGeckoResponse struct {
	Bitcoin struct {
		UAH int `json:"uah"`
	} `json:"bitcoin"`
}

func (repository CoinGeckoRepository) GetRate(from, to string) (float64, error) {
	var response coinGeckoResponse

	resp, err := repository.client.R().
		SetQueryParam("ids", supportedRatesKeys[strings.ToLower(from)]).
		SetQueryParam("vs_currencies", supportedRatesKeys[strings.ToLower(to)]).
		Get(repository.config.URL)

	if err != nil {
		return 0, fmt.Errorf("failed to perform API request: %w", err)
	}

	logrus.Infof("CoinGecko API response: %s", resp.String())

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("unexpected API response: %s", resp.Status())
	}

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return 0, fmt.Errorf("failed to parse API response: %w", err)
	}

	price := float64(response.Bitcoin.UAH)

	return price, nil
}
