package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationService interface {
	Send([]string, float64) error
}

type NotificationHandler struct {
	externalService     RateService
	notificationService NotificationService
	storage             StorageRepository
	logger              logger
}

func NewNotificationHandler(externalService RateService, notificationService NotificationService,
	storage StorageRepository, logger logger) *NotificationHandler {
	return &NotificationHandler{
		externalService:     externalService,
		notificationService: notificationService,
		storage:             storage,
		logger:              logger,
	}
}

func (handler *NotificationHandler) SendEmails(c *gin.Context) {
	emails, err := handler.storage.AllEmails()

	if err != nil {
		handler.logger.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	rate, err := handler.externalService.CurrentRate("BTC", "UAH")

	if err != nil {
		handler.logger.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	err = handler.notificationService.Send(emails, rate)

	if err != nil {
		handler.logger.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
