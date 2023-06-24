package service_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"ses.genesis.com/exchange-web-service/src/main/service"
)

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

func TestGetRate(t *testing.T) {
	t.Run("shouldGetRate", func(t *testing.T) {
		externalService := &MockExternalService{}
		internalService := service.NewMainService(externalService, nil, nil)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		internalService.GetRate(ctx)
		if ctx.Writer.Status() != http.StatusOK {
			t.Error("Status code is not 200")
		}
	})

	t.Run("shouldNotGetRate", func(t *testing.T) {
		externalService := &MockExternalServiceFail{}
		internalService := service.NewMainService(externalService, nil, nil)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		internalService.GetRate(ctx)
		if ctx.Writer.Status() != http.StatusInternalServerError {
			t.Errorf("Status code is not 500")
		}
	})
}
