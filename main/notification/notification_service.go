package notification

import "net/smtp"

type NotifyProtocolService interface {
	Authenticate(config AuthConfig) smtp.Auth
	SendMessage(auth smtp.Auth, config AuthConfig, to []string, massage []byte) error
}
