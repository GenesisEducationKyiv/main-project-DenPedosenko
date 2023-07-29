package cmd

import (
	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
)

type MainService interface {
	GetRate(c *gin.Context)
	PostEmail(c *gin.Context)
	GetEmails(c *gin.Context)
	SendEmails(c *gin.Context)
}

type SagaService interface {
	PostSubscribe(c *gin.Context)
	PostSaveEmail(c *gin.Context) any
	PostRemoveEmail(c *gin.Context) any
}

type GinRouter struct {
	client  *gin.Engine
	service MainService
	saga    SagaService
}

func NewGinRouter(service MainService, sagaService SagaService) *GinRouter {
	return &GinRouter{
		client:  gin.Default(),
		service: service,
		saga:    sagaService,
	}
}

func (router GinRouter) CreateRoutes() *gin.Engine {
	router.client.GET("api/rate", router.service.GetRate)
	router.client.GET("api/subscribe", router.service.GetEmails)
	router.client.POST("api/subscribe", router.service.PostEmail)
	router.client.POST("api/subscribe_saga", router.saga.PostSubscribe)
	router.client.POST("save_email", dtmutil.WrapHandler2(router.saga.PostSaveEmail))
	router.client.POST("remove_email", dtmutil.WrapHandler2(router.saga.PostRemoveEmail))
	router.client.POST("api/sendEmails", router.service.SendEmails)

	return router.client
}
