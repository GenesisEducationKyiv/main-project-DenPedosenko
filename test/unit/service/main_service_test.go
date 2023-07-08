package testservice

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"ses.genesis.com/exchange-web-service/main/presentation/handler"
	"ses.genesis.com/exchange-web-service/main/presentation/handler/errormapper"
	"ses.genesis.com/exchange-web-service/main/service"

	"github.com/gin-gonic/gin"
)

func TestServiceError(t *testing.T) {
	externalService := &MockExternalServiceFail{}
	persistentService := &MockPersistentRepository{}
	notificationService := &MockNotificationServiceFail{}
	mapper := errormapper.NewStorageErrorToHTTPMapper()

	rateController := handler.NewRateHandler(externalService)
	emailController := handler.NewEmailHandler(persistentService, mapper)
	notificationController := handler.NewNotificationHandler(externalService, notificationService, persistentService)

	t.Run("shouldNotGetRate", func(t *testing.T) {
		internalService := service.NewMainService(rateController, emailController, notificationController)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		internalService.GetRate(ctx)
		if ctx.Writer.Status() != http.StatusInternalServerError {
			t.Errorf("Status code is not 500")
		}
	})

	t.Run("shouldNotPostEmail", func(t *testing.T) {
		persistentService := &MockPersistentServiceFail{}
		emailControllerFail := handler.NewEmailHandler(persistentService, mapper)

		internalService := service.NewMainService(rateController, emailControllerFail, notificationController)
		ctx := getTestRequestContext()
		internalService.PostEmail(ctx)
		if ctx.Writer.Status() != http.StatusInternalServerError {
			t.Errorf("Status code is not 500")
		}
	})

	t.Run("shouldNotGetEmails", func(t *testing.T) {
		persistentService := &MockPersistentServiceFail{}
		emailControllerFail := handler.NewEmailHandler(persistentService, mapper)
		internalService := service.NewMainService(rateController, emailControllerFail, notificationController)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		internalService.GetEmails(ctx)
		if ctx.Writer.Status() != http.StatusInternalServerError {
			t.Errorf("Status code is not 500")
		}
	})

	t.Run("shouldNotSendEmails", func(t *testing.T) {
		internalService := service.NewMainService(rateController, emailController, notificationController)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		internalService.SendEmails(ctx)
		if ctx.Writer.Status() != http.StatusInternalServerError {
			t.Errorf("Status code is not 500")
		}
	})
}

func TestServiceSuccess(t *testing.T) {
	externalService := &MockExternalService{}
	persistentService := &MockPersistentRepository{}
	notificationService := &MockNotificationService{}
	mapper := errormapper.NewStorageErrorToHTTPMapper()

	rateController := handler.NewRateHandler(externalService)
	emailController := handler.NewEmailHandler(persistentService, mapper)
	notificationController := handler.NewNotificationHandler(externalService, notificationService, persistentService)

	internalService := service.NewMainService(rateController, emailController, notificationController)

	t.Run("shouldGetRate", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		internalService.GetRate(ctx)
		if ctx.Writer.Status() != http.StatusOK {
			t.Error("Status code is not 200")
		}
	})

	t.Run("shouldSendEmails", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		internalService.SendEmails(ctx)
		if ctx.Writer.Status() != http.StatusOK {
			t.Errorf("Status code is not 200")
		}
	})

	t.Run("shouldGetEmails", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		internalService.GetEmails(ctx)
		if ctx.Writer.Status() != http.StatusOK {
			t.Errorf("Status code is not 200")
		}
	})

	t.Run("shouldPostEmail", func(t *testing.T) {
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
}

func getTestRequestContext() *gin.Context {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	buf := new(bytes.Buffer)
	buf.WriteString("email=test@gmail.com")

	ctx.Request, _ = http.NewRequest("POST", "/api/subscribe", buf)
	ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return ctx
}
