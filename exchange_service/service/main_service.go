package service

import (
	"exchange-web-service/presentation/handler"

	"github.com/gin-gonic/gin"
)

type Rate interface {
	GetRate(c *gin.Context)
}

type Email interface {
	PostEmail(c *gin.Context)
	GetEmails(c *gin.Context)
}

type Notification interface {
	SendEmails(c *gin.Context)
}

type MainService struct {
	rate         Rate
	email        Email
	notification Notification
}

func (service *MainService) GetRate(c *gin.Context) {
	service.rate.GetRate(c)
}

func (service *MainService) PostEmail(c *gin.Context) {
	service.email.PostEmail(c)
}

func (service *MainService) GetEmails(c *gin.Context) {
	service.email.GetEmails(c)
}

func (service *MainService) SendEmails(c *gin.Context) {
	service.notification.SendEmails(c)
}

func NewMainService(rate *handler.RateHandler, email *handler.EmailHandler,
	notification *handler.NotificationHandler) *MainService {
	return &MainService{
		rate:         rate,
		email:        email,
		notification: notification,
	}
}
