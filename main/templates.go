package main

import _ "embed"

//go:embed template/message.html
var message string

func getMessageTemplate() string {
	return message
}
