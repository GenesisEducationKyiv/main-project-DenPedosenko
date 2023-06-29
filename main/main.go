package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	var service = initialize()

	router := gin.Default()
	router.GET("api/rate", service.GetRate)
	router.GET("api/subscribe", service.GetEmails)
	router.POST("api/subscribe", service.PostEmail)
	router.POST("api/sendEmails", service.SendEmails)
	err := router.Run("localhost:8080")

	if err != nil {
		log.Fatal(err)
	}
}
