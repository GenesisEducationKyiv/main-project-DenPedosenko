package service

import "fmt"

type MockExternalService struct {
}

func (s *MockExternalService) GetCurrentBTCToUAHRate() (float64, error) {
	return 1, nil
}

type MockExternalServiceFail struct {
}

func (s *MockExternalServiceFail) GetCurrentBTCToUAHRate() (float64, error) {
	return -1, fmt.Errorf("failed to get rate")
}
