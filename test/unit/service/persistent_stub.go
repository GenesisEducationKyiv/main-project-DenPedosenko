package testservice

import (
	"errors"
	"fmt"

	"ses.genesis.com/exchange-web-service/main/persistent"
)

type MockPersistentService struct {
	emails []string
}

func (s *MockPersistentService) AllEmails() ([]string, error) {
	return s.emails, nil
}

func (s *MockPersistentService) SaveEmailToStorage(email string) persistent.StorageError {
	if s.IsEmailAlreadyExists(email) {
		return persistent.StorageError{
			Code: persistent.Conflict,
			Err:  errors.New("email already exists"),
		}
	}

	s.emails = append(s.emails, email)

	return persistent.StorageError{
		Code: -1,
		Err:  nil,
	}
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

func (s *MockPersistentServiceFail) SaveEmailToStorage(_ string) persistent.StorageError {
	return persistent.StorageError{Err: fmt.Errorf("failed to save email"), Code: persistent.UnknownError}
}

func (s *MockPersistentServiceFail) IsEmailAlreadyExists(_ string) bool {
	return false
}
