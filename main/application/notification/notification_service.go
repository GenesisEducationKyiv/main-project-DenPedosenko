package notification

import (
	"net/smtp"
)

type AuthConfig struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func NewAuthConfig(from, password, smtpHost, smtpPort string) AuthConfig {
	return AuthConfig{
		from:     from,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

type NotifyProtocolService interface {
	Authenticate(config AuthConfig) smtp.Auth
	SendMessage(auth smtp.Auth, config AuthConfig, to []string, massage []byte) error
}
