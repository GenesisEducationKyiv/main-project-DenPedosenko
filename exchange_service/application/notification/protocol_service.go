package notification

import (
	"net/smtp"
)

type SMTPProtocolService struct {
}

func NewSMTPProtocolService() NotifyProtocolService {
	return &SMTPProtocolService{}
}

func (protocol *SMTPProtocolService) Authenticate(config AuthConfig) smtp.Auth {
	return smtp.PlainAuth("", config.from, config.password, config.smtpHost)
}

func (protocol *SMTPProtocolService) SendMessage(auth smtp.Auth, config AuthConfig, to []string, massage []byte) error {
	return smtp.SendMail(config.smtpHost+":"+config.smtpPort, auth, config.from, to, massage)
}
