package main

import (
	"log"
	"ses.genesis.com/exchange-web-service/main/cmd"
)

func main() {
	err := cmd.Router().Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
}
