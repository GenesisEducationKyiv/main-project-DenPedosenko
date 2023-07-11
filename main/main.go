package main

import (
	"log"

	"ses.genesis.com/exchange-web-service/main/cmd"
)

func main() {
	app := cmd.NewApplication()
	err := app.Router.CreateRoutes().Run("localhost:8080")

	if err != nil {
		log.Fatal(err)
	}
}
