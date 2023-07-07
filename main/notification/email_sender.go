package notification

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"text/template"

	"ses.genesis.com/exchange-web-service/main/config"
)

const mimeHeaders = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

type EmailSender struct {
	ctx      context.Context
	protocol NotifyProtocolService
	template *template.Template
}

type AuthConfig struct {
	from     string
	password string
	smtpHost string
	smtpPort string
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

	var authConfig = AuthConfig{
		from:     conf.EmailUser,
		password: conf.EmailPassword,
		smtpHost: conf.EmailHost,
		smtpPort: conf.EmailPort,
	}

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

	body.Write([]byte(fmt.Sprintf("Subject: Current BTC to UAH exchange rate \n%s\n\n", mimeHeaders)))

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
