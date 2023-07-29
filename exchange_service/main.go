package main

import (
	"exchange-web-service/cmd"
	"log"
)

func main() {
	app := cmd.NewApplication()
	err := app.Router.CreateRoutes().Run(":8080")

	if err != nil {
		log.Fatal(err)
	}

	app.Logger.Close()
}
