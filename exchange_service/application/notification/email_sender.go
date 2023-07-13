package notification

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/smtp"
	"text/template"

	"exchange-web-service/domain/config"
)

const mimeHeaders = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

type NotifyProtocolService interface {
	Authenticate(config AuthConfig) smtp.Auth
	SendMessage(auth smtp.Auth, config AuthConfig, to []string, massage []byte) error
}

type EmailSender struct {
	ctx      context.Context
	protocol NotifyProtocolService
	template *template.Template
}

func NewEmailSender(ctx context.Context, protocol NotifyProtocolService) *EmailSender {
	t, errTemplate := template.New("message").Parse(getMessageTemplate())
	if errTemplate != nil {
		log.Fatal(errTemplate)
	}

	return &EmailSender{
		ctx:      ctx,
		template: t,
		protocol: protocol,
	}
}

func (sender *EmailSender) Send(to []string, rate float64) error {
	conf := config.GetConfigFromContext(sender.ctx)

	var authConfig = NewAuthConfig(conf.EmailUser, conf.EmailPassword, conf.EmailHost, conf.EmailPort)

	auth := sender.protocol.Authenticate(authConfig)

	body, errBody := sender.getMessageBody(rate)
	if errBody != nil {
		log.Fatal(errBody)
		return errBody
	}

	errSendMail := sender.protocol.SendMessage(auth, authConfig, to, body.Bytes())
	if errSendMail != nil {
		return errSendMail
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
		log.Fatal(err)
		return nil, err
	}

	return &body, nil
}
