package service

import (
	"fmt"
	"ses.genesis.com/exchange-web-service/src/main/persistent"
)

type MockPersistentService struct {
	emails []string
}

func (s *MockPersistentService) AllEmails() ([]string, error) {
	return s.emails, nil
}

func (s *MockPersistentService) SaveEmailToStorage(email string) (int, error) {
	if s.IsEmailAlreadyExists(email) {
		return int(persistent.Conflict), nil
	}

	s.emails = append(s.emails, email)

	return int(persistent.OK), nil
}

func (s *MockPersistentService) IsEmailAlreadyExists(email string) bool {
	for _, e := range s.emails {
		if e == email {
			return true
		}
	}

	return false
}

type MockPersistentServiceFail struct {
}

func (s *MockPersistentServiceFail) AllEmails() ([]string, error) {
	return nil, fmt.Errorf("failed to get emails")
}

func (s *MockPersistentServiceFail) SaveEmailToStorage(_ string) (int, error) {
	return int(persistent.UnknownError), fmt.Errorf("failed to save email")
}

func (s *MockPersistentServiceFail) IsEmailAlreadyExists(_ string) bool {
	return false
}
