package main

import (
	"log"
	"log_consumer/consumer"
)

func main() {
	logConsumer, err := consumer.NewLoggerConsumer()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Logger Consumer started. Waiting for logs...")

	err = logConsumer.ConsumeLogs()
	if err != nil {
		log.Fatal(err)
	}
}
