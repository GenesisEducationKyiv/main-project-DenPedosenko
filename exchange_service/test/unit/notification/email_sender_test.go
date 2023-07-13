package notification_test

import (
	"net/smtp"
	"testing"

	"exchange-web-service/application/notification"

	"exchange-web-service/domain/config"
)

type MockSMTPProtocolService struct {
	message string
}

func (protocol *MockSMTPProtocolService) Authenticate(_ notification.AuthConfig) smtp.Auth {
	return nil
}

func (protocol *MockSMTPProtocolService) SendMessage(_ smtp.Auth, _ notification.AuthConfig, _ []string, massage []byte) error {
	protocol.message = string(massage)
	return nil
}

func TestSend(t *testing.T) {
	configLoader := config.NewConfigLoader("../../application.yaml")
	ctx, _ := configLoader.GetContext()

	t.Run("shouldSendEmail", func(t *testing.T) {
		protocol := &MockSMTPProtocolService{}
		sender := notification.NewEmailSender(ctx, protocol)
		err := sender.Send([]string{"test@gmail.com"}, 1.0)
		if err != nil {
			t.Errorf("Err: %s", err)
		}
		if protocol.message == "" {
			t.Errorf("Message is empty")
		}
	})
}
