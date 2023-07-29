package testservice

import "log"

type TestLogger struct {
}

func (t TestLogger) Info(msg string) {
	log.Print(msg)
}

func (t TestLogger) Error(msg string) {
	log.Print(msg)
}

func (t TestLogger) Debug(msg string) {
	log.Print(msg)
}

func (t TestLogger) Close() {
	log.Print("Closed")
}
