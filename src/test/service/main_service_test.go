package service

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"ses.genesis.com/exchange-web-service/src/main/service"
)

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

func TestPostEmail(t *testing.T) {
	t.Run("shouldPostEmail", func(t *testing.T) {
		persistentService := &MockPersistentService{}
		internalService := service.NewMainService(nil, persistentService, nil)
		ctx := getTestRequestContext()
		internalService.PostEmail(ctx)
		if ctx.Writer.Status() != http.StatusOK {
			t.Errorf("Status code is not 200")
		}

		if persistentService.emails[0] != "test@gmail.com" {
			t.Errorf("Email is not saved")
		}
	})

	t.Run("shouldNotPostEmailWithConflict", func(t *testing.T) {
		persistentService := &MockPersistentService{}
		internalService := service.NewMainService(nil, persistentService, nil)
		persistentService.emails = append(persistentService.emails, "test@gmail.com")
		ctx := getTestRequestContext()
		internalService.PostEmail(ctx)
		if ctx.Writer.Status() != http.StatusConflict {
			t.Errorf("Status code is not 409")
		}

		if persistentService.emails[0] != "test@gmail.com" {
			t.Errorf("Email is not saved")
		}
	})

	t.Run("shouldNotPostEmail", func(t *testing.T) {
		persistentService := &MockPersistentServiceFail{}
		internalService := service.NewMainService(nil, persistentService, nil)
		ctx := getTestRequestContext()
		internalService.PostEmail(ctx)
		if ctx.Writer.Status() != http.StatusInternalServerError {
			t.Errorf("Status code is not 500")
		}
	})
}

func TestGetEmails(t *testing.T) {
	t.Run("shouldGetEmails", func(t *testing.T) {
		persistentService := &MockPersistentService{}
		internalService := service.NewMainService(nil, persistentService, nil)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		internalService.GetEmails(ctx)
		if ctx.Writer.Status() != http.StatusOK {
			t.Errorf("Status code is not 200")
		}
	})

	t.Run("shouldNotGetEmails", func(t *testing.T) {
		persistentService := &MockPersistentServiceFail{}
		internalService := service.NewMainService(nil, persistentService, nil)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		internalService.GetEmails(ctx)
		if ctx.Writer.Status() != http.StatusInternalServerError {
			t.Errorf("Status code is not 500")
		}
	})
}

func TestSendEmails(t *testing.T) {
	t.Run("shouldSendEmails", func(t *testing.T) {
		externalService := &MockExternalService{}
		persistentService := &MockPersistentService{}
		notificationService := &MockNotificationService{}
		internalService := service.NewMainService(externalService, persistentService, notificationService)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		internalService.SendEmails(ctx)
		if ctx.Writer.Status() != http.StatusOK {
			t.Errorf("Status code is not 200")
		}
	})

	t.Run("shouldNotSendEmails", func(t *testing.T) {
		externalService := &MockExternalServiceFail{}
		persistentService := &MockPersistentService{}
		notificationService := &MockNotificationServiceFail{}
		internalService := service.NewMainService(externalService, persistentService, notificationService)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		internalService.SendEmails(ctx)
		if ctx.Writer.Status() != http.StatusInternalServerError {
			t.Errorf("Status code is not 500")
		}
	})

}

func getTestRequestContext() *gin.Context {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	buf := new(bytes.Buffer)
	buf.WriteString("email=test@gmail.com")

	ctx.Request, _ = http.NewRequest("POST", "/api/subscribe", buf)
	ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return ctx
}
