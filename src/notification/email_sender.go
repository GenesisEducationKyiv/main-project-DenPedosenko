package notification

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/smtp"
	"text/template"

	"ses.genesis.com/exchange-web-service/src/config"
)

type EmailSender struct {
	ctx context.Context
}

func NewEmailSender(ctx context.Context) NotificationService {
	return &EmailSender{
		ctx: ctx,
	}
}

func (sender *EmailSender) Send(to []string, rate float64) error {
	conf := config.GetConfig(sender.ctx)
	from := conf.EmailUser
	password := conf.EmailPassword
	smtpHost := conf.EmailHost
	smtpPort := conf.EmailPort

	auth := smtp.PlainAuth("", from, password, smtpHost)

	body, errBody := getMessageBody(rate)
	if errBody != nil {
		log.Fatal(errBody)
		return errBody
	}

	errSendMail := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if errSendMail != nil {
		return errSendMail
	}

	return nil
}

func getMessageBody(rate float64) (*bytes.Buffer, error) {
	t, errTemplate := template.New("message").Parse(getMessageTemplate())
	if errTemplate != nil {
		log.Fatal(errTemplate)
		return nil, errTemplate
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Current BTC to UAH exchange rate \n%s\n\n", mimeHeaders)))

	err := t.Execute(&body, struct {
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
