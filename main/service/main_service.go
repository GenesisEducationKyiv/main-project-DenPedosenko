package service

import (
	"github.com/gin-gonic/gin"
)

type InternalService interface {
	GetRate(*gin.Context)
	PostEmail(*gin.Context)
	GetEmails(*gin.Context)
	SendEmails(*gin.Context)
}

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
	rate         *RateController
	email        *EmailController
	notification *NotificationController
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

func NewMainService(rate *RateController, email *EmailController, notification *NotificationController) *MainService {
	return &MainService{
		rate:         rate,
		email:        email,
		notification: notification,
	}
}
