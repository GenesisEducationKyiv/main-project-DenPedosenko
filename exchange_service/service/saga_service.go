package service

import (
	"exchange-web-service/persistent"
	"exchange-web-service/presentation/handler"
	"log"
	"net/http"
	"os"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/gin-gonic/gin"
)

type config struct {
	dtmCoordinatorAddress string
	exchangeServerURL     string
	customersServerURL    string
}
type SagaService struct {
	storageRepository  handler.StorageRepository
	storageErrorMapper handler.StorageErrorMapper[persistent.StorageError, int]
	config             config
}

type EmailRequest struct {
	Email string `json:"email"`
}

func NewEmailRequest(email string) *EmailRequest {
	return &EmailRequest{
		Email: email,
	}
}

func NewSagaService(storageRepository handler.StorageRepository,
	storageErrorMapper handler.StorageErrorMapper[persistent.StorageError, int]) *SagaService {
	return &SagaService{
		storageRepository:  storageRepository,
		storageErrorMapper: storageErrorMapper,
		config: config{
			dtmCoordinatorAddress: os.Getenv("DTM_COORDINATOR"),
			exchangeServerURL:     os.Getenv("EXCHANGE_SERVICE_URL"),
			customersServerURL:    os.Getenv("CUSTOMERS_SERVICE_URL"),
		},
	}
}

func (service *SagaService) PostSubscribe(c *gin.Context) {
	request := c.Request
	writer := c.Writer
	err := request.ParseForm()
	headerContentType := request.Header.Get("Content-Type")

	if headerContentType != "application/x-www-form-urlencoded" {
		writer.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	req := NewEmailRequest(request.FormValue("email"))
	log.Printf("req: %v", req)

	globalID := dtmcli.MustGenGid(service.config.dtmCoordinatorAddress)

	err = dtmcli.
		NewSaga(service.config.dtmCoordinatorAddress, globalID).
		Add(service.config.exchangeServerURL+"/save_email", service.config.exchangeServerURL+"/remove_email", req).
		Add(service.config.customersServerURL+"/create", "", nil).
		Submit()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}

	writer.WriteHeader(http.StatusOK)
}

func (service *SagaService) PostSaveEmail(c *gin.Context) any {
	processEmail(c, service, service.storageRepository.Save)

	data := map[string]string{
		"status":  "success",
		"message": "Email saved successfully",
	}

	return data
}

func (service *SagaService) PostRemoveEmail(c *gin.Context) any {
	processEmail(c, service, service.storageRepository.Remove)

	data := map[string]string{
		"status":  "success",
		"message": "Email removed successfully",
	}

	return data
}

func processEmail(c *gin.Context, service *SagaService, save func(string) persistent.StorageError) {
	resp := &EmailRequest{}
	request := c.Request
	writer := c.Writer
	err := request.ParseForm()

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.BindJSON(resp)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	errSave := save(resp.Email)

	if errSave.Err != nil {
		writer.WriteHeader(service.storageErrorMapper.MapError(errSave))
		return
	}

	writer.WriteHeader(http.StatusOK)
}
