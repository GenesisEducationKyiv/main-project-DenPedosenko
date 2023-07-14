package cmd

import (
	"github.com/gin-gonic/gin"
)

type MainService interface {
	GetRate(c *gin.Context)
	PostEmail(c *gin.Context)
	GetEmails(c *gin.Context)
	SendEmails(c *gin.Context)
}

type GinRouter struct {
	client  *gin.Engine
	service MainService
}

func NewGinRouter(service MainService) *GinRouter {
	return &GinRouter{
		client:  gin.Default(),
		service: service,
	}
}

func (router GinRouter) CreateRoutes() *gin.Engine {
	router.client.GET("api/rate", router.service.GetRate)
	router.client.GET("api/subscribe", router.service.GetEmails)
	router.client.POST("api/subscribe", router.service.PostEmail)
	router.client.POST("api/sendEmails", router.service.SendEmails)

	return router.client
}
