package notification

import (
	"bytes"
	"context"
	"exchange-web-service/domain/config"
	"fmt"
	"net/smtp"
	"text/template"
)

const mimeHeaders = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

type NotifyProtocolService interface {
	Authenticate(config AuthConfig) smtp.Auth
	SendMessage(auth smtp.Auth, config AuthConfig, to []string, massage []byte) error
}

type logger interface {
	Info(string string)
	Error(string string)
	Debug(string string)
	Close()
}

type EmailSender struct {
	ctx      context.Context
	protocol NotifyProtocolService
	template *template.Template
	logger   logger
}

func NewEmailSender(ctx context.Context, protocol NotifyProtocolService, logger logger) *EmailSender {
	t, err := template.New("message").Parse(getMessageTemplate())
	if err != nil {
		logger.Error(err.Error())
	}

	return &EmailSender{
		ctx:      ctx,
		template: t,
		protocol: protocol,
		logger:   logger,
	}
}

func (sender *EmailSender) Send(to []string, rate float64) error {
	conf := config.GetConfigFromContext(sender.ctx)

	var authConfig = NewAuthConfig(conf.EmailUser, conf.EmailPassword, conf.EmailHost, conf.EmailPort)

	auth := sender.protocol.Authenticate(authConfig)

	body, err := sender.getMessageBody(rate)
	if err != nil {
		sender.logger.Error(err.Error())
		return err
	}

	err = sender.protocol.SendMessage(auth, authConfig, to, body.Bytes())
	if err != nil {
		sender.logger.Error(err.Error())
		return err
	}

	return nil
}

func (sender *EmailSender) getMessageBody(rate float64) (*bytes.Buffer, error) {
	var body bytes.Buffer

	body.Write([]byte(fmt.Sprintf("Subject: Current BTC to UAH exchange_service rate \n%s\n\n", mimeHeaders)))

	err := sender.template.Execute(&body, struct {
		Rate string
	}{
		Rate: fmt.Sprintf("%f", rate),
	})

	if err != nil {
		sender.logger.Error(err.Error())
		return nil, err
	}

	return &body, nil
}
