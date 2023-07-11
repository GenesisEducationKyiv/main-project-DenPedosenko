package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RateHandler struct {
	externalService RateService
}

type RateService interface {
	CurrentRate(from string, to string) (float64, error)
}

func NewRateHandler(externalService RateService) *RateHandler {
	return &RateHandler{
		externalService: externalService,
	}
}

func (handler RateHandler) GetRate(c *gin.Context) {
	rate, err := handler.externalService.CurrentRate("BTC", "UAH")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, rate)
}
