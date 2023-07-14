package notification_test

import (
	notification2 "exchange-web-service/main/application/notification"
	"exchange-web-service/main/domain/config"
	"net/smtp"
	"testing"
)

type MockSMTPProtocolService struct {
	message string
}

func (protocol *MockSMTPProtocolService) Authenticate(_ notification2.AuthConfig) smtp.Auth {
	return nil
}

func (protocol *MockSMTPProtocolService) SendMessage(_ smtp.Auth, _ notification2.AuthConfig, _ []string, massage []byte) error {
	protocol.message = string(massage)
	return nil
}

func TestSend(t *testing.T) {
	configLoader := config.NewConfigLoader("../../application.yaml")
	ctx, _ := configLoader.GetContext()

	t.Run("shouldSendEmail", func(t *testing.T) {
		protocol := &MockSMTPProtocolService{}
		sender := notification2.NewEmailSender(ctx, protocol)
		err := sender.Send([]string{"test@gmail.com"}, 1.0)
		if err != nil {
			t.Errorf("Err: %s", err)
		}
		if protocol.message == "" {
			t.Errorf("Message is empty")
		}
	})
}
