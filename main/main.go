package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("api/rate", getRate)
	router.GET("api/subscribe", getEmails)
	router.POST("api/subscribe", postEmail)
	router.POST("api/sendEmails", sendEmails)
	err := router.Run("localhost:8080")

	if err != nil {
		log.Fatal(err)
	}
}
