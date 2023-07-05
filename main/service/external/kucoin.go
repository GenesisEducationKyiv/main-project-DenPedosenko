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

type KuCoinAPIProvider struct {
	config *config.ConfigAPI
	client *resty.Client
}

type KuCoinResponse struct {
	Code string `json:"code"`
	Data map[string]string
}

func NewKuCoinRepository(conf *config.ConfigAPI, client *resty.Client) *KuCoinAPIProvider {
	return &KuCoinAPIProvider{
		config: conf,
		client: client,
	}
}

func (repository KuCoinAPIProvider) GetRate(from, to string) (float64, error) {
	var response KuCoinResponse

	resp, err := repository.client.R().
		SetQueryParam("currencies", from).
		SetQueryParam("base", to).
		Get(repository.config.URL)

	if err != nil {
		return 0, fmt.Errorf("failed to perform API request: %w", err)
	}

	logrus.Infof("KuCoin API Data: %s", resp.String())

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("unexpected API response: %s", resp.Status())
	}

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return 0, fmt.Errorf("failed to parse API Data: %w", err)
	}

	price, err := strconv.ParseFloat(response.Data[from], 64)

	if err != nil {
		return 0, fmt.Errorf("failed to parse API Data: %w", err)
	}

	return price, nil
}
