package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RateHandler struct {
	externalService RateService
	logger          logger
}

type RateService interface {
	CurrentRate(from string, to string) (float64, error)
}

func NewRateHandler(externalService RateService, logger logger) *RateHandler {
	return &RateHandler{
		externalService: externalService,
		logger:          logger,
	}
}

func (handler RateHandler) GetRate(c *gin.Context) {
	rate, err := handler.externalService.CurrentRate("BTC", "UAH")
	if err != nil {
		handler.logger.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	handler.logger.Info(fmt.Sprintf("Rate is %f", rate))
	c.IndentedJSON(http.StatusOK, rate)
}
