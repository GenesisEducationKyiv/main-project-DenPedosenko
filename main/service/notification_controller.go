package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationService interface {
	Send([]string, float64) error
}

type NotificationController struct {
	externalService     RateService
	notificationService NotificationService
	storage             StorageRepository
}

func NewNotificationController(externalService RateService, notificationService NotificationService,
	storage StorageRepository) *NotificationController {
	return &NotificationController{
		externalService:     externalService,
		notificationService: notificationService,
		storage:             storage,
	}
}

func (us *NotificationController) SendEmails(c *gin.Context) {
	emails, err := us.storage.AllEmails()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rate, err := us.externalService.CurrentRate("BTC", "UAH")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = us.notificationService.Send(emails, rate)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
