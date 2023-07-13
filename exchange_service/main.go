package main

import (
	"log"

	"exchange-web-service/cmd"
)

func main() {
	app := cmd.NewApplication()
	err := app.Router.CreateRoutes().Run("localhost:8080")

	if err != nil {
		log.Fatal(err)
	}
}
