package service

import "fmt"

type MockNotificationService struct {
}

func (s *MockNotificationService) Send(_ []string, _ float64) error {
	return nil
}

type MockNotificationServiceFail struct {
}

func (s *MockNotificationServiceFail) Send(_ []string, _ float64) error {
	return fmt.Errorf("failed to send email")
}
