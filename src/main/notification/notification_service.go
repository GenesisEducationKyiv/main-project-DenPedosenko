package notification

import "net/smtp"

type NotifyService interface {
	Send([]string, float64) error
}

type NotifyProtocolService interface {
	Authenticate(config AuthConfig) smtp.Auth
	SendMessage(auth smtp.Auth, config AuthConfig, to []string, massage []byte) error
}
