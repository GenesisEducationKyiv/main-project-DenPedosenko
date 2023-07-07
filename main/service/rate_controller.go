package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RateController struct {
	externalService RateService
}

type RateService interface {
	CurrentRate(from string, to string) (float64, error)
}

func NewRateController(externalService RateService) *RateController {
	return &RateController{
		externalService: externalService,
	}
}

func (controller RateController) GetRate(c *gin.Context) {
	rate, err := controller.externalService.CurrentRate("BTC", "UAH")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, rate)
}
