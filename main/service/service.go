package service

import "github.com/gin-gonic/gin"

type InternalService interface {
	GetRate(*gin.Context)
	PostEmail(*gin.Context)
	GetEmails(*gin.Context)
	SendEmails(*gin.Context)
}
