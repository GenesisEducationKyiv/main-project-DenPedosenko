package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"ses.genesis.com/exchange-web-service/main/config"
)

type KuCoinAPIRepository struct {
	config *config.ConfigAPI
	client *resty.Client
}

type kuCoinResponse struct {
	Code string `json:"code"`
	Data struct {
		BTC string `json:"BTC"`
	} `json:"data"`
}

func NewKuCoinRepository(conf *config.ConfigAPI, client *resty.Client) *KuCoinAPIRepository {
	return &KuCoinAPIRepository{
		config: conf,
		client: client,
	}
}

func (repository KuCoinAPIRepository) GetRate(from, to string) (float64, error) {
	var response kuCoinResponse

	resp, err := repository.client.R().
		SetQueryParam("currencies", from).
		SetQueryParam("base", to).
		Get(repository.config.URL)

	if err != nil {
		return 0, fmt.Errorf("failed to perform API request: %w", err)
	}

	logrus.Infof("KuCoin API response: %s", resp.String())

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("unexpected API response: %s", resp.Status())
	}

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return 0, fmt.Errorf("failed to parse API response: %w", err)
	}

	price, err := strconv.ParseFloat(response.Data.BTC, 64)

	if err != nil {
		return 0, fmt.Errorf("failed to parse API response: %w", err)
	}

	return price, nil
}
