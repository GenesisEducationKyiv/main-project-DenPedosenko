package notification

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"
)

type EmailSender struct {
}

func NewEmailSender() NotificationService {
	return &EmailSender{}
}

func (sender *EmailSender) Send(to []string, rate float64) error {
	from := "test.sender.genesis.ses@gmail.com"
	password := "uotikrysrztcgdiq" //nolint:gosec

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

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

	_ = t.Execute(&body, struct {
		Rate string
	}{
		Rate: fmt.Sprintf("%f", rate),
	})

	return &body, nil
}
