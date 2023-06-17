package main

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"
)

func send(to []string) error {
	from := "test.sender.genesis.ses@gmail.com"
	password := "uotikrysrztcgdiq" //nolint:gosec

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, errTemplate := template.New("message").Parse(getMessageTemplate())
	if errTemplate != nil {
		log.Fatal(errTemplate)
		return errTemplate
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Current BTC to UAH exchange rate \n%s\n\n", mimeHeaders)))

	rate, errRateExtr := getCurrentBTCToUAHRate()
	if errRateExtr != nil {
		log.Fatal(errRateExtr)
	}

	_ = t.Execute(&body, struct {
		Rate string
	}{
		Rate: fmt.Sprintf("%f", rate),
	})

	errSendMail := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if errSendMail != nil {
		return errSendMail
	}

	return nil
}
