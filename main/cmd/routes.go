package cmd

import "github.com/gin-gonic/gin"

func Router() *gin.Engine {
	var service = initialize()

	router := gin.Default()
	router.GET("api/rate", service.GetRate)
	router.GET("api/subscribe", service.GetEmails)
	router.POST("api/subscribe", service.PostEmail)
	router.POST("api/sendEmails", service.SendEmails)

	return router
}
