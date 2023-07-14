package testservice

import "fmt"

type MockExternalService struct {
}

func (s *MockExternalService) CurrentRate(_, _ string) (float64, error) {
	return 1, nil
}

type MockExternalServiceFail struct {
}

func (s *MockExternalServiceFail) CurrentRate(_, _ string) (float64, error) {
	return -1, fmt.Errorf("failed to get rate")
}
