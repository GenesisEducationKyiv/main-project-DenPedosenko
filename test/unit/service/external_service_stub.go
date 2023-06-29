package testservice

import "fmt"

type MockExternalService struct {
}

func (s *MockExternalService) CurrentBTCToUAHRate() (float64, error) {
	return 1, nil
}

type MockExternalServiceFail struct {
}

func (s *MockExternalServiceFail) CurrentBTCToUAHRate() (float64, error) {
	return -1, fmt.Errorf("failed to get rate")
}
