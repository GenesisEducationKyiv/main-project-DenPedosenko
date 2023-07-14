package main

import (
	"exchange-web-service/main/cmd"
	"log"
)

func main() {
	app := cmd.NewApplication()
	err := app.Router.CreateRoutes().Run("localhost:8080")

	if err != nil {
		log.Fatal(err)
	}
}
