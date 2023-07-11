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
}

func NewNotificationHandler(externalService RateService, notificationService NotificationService,
	storage StorageRepository) *NotificationHandler {
	return &NotificationHandler{
		externalService:     externalService,
		notificationService: notificationService,
		storage:             storage,
	}
}

func (handler *NotificationHandler) SendEmails(c *gin.Context) {
	emails, err := handler.storage.AllEmails()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rate, err := handler.externalService.CurrentRate("BTC", "UAH")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = handler.notificationService.Send(emails, rate)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
