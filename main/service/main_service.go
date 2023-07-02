package service

import (
	"net/http"

	"ses.genesis.com/exchange-web-service/main/persistent"
	"ses.genesis.com/exchange-web-service/main/service/external"

	"ses.genesis.com/exchange-web-service/main/notification"
	"ses.genesis.com/exchange-web-service/main/service/errormapper"

	"github.com/gin-gonic/gin"
)

type InternalService interface {
	GetRate(*gin.Context)
	PostEmail(*gin.Context)
	GetEmails(*gin.Context)
	SendEmails(*gin.Context)
}

type MainService struct {
	externalService     external.RateService
	notificationService notification.NotifyService
	storageRepository   persistent.StorageRepository
	storageErrorMapper  errormapper.StorageErrorMapper[persistent.StorageError, int]
}

func NewMainService(externalService external.RateService, persistentService persistent.StorageRepository,
	notificationService notification.NotifyService,
	storageErrorMapper errormapper.StorageErrorMapper[persistent.StorageError, int]) *MainService {
	return &MainService{
		externalService:     externalService,
		storageRepository:   persistentService,
		notificationService: notificationService,
		storageErrorMapper:  storageErrorMapper}
}

func (service *MainService) GetRate(c *gin.Context) {
	rate, err := service.externalService.CurrentRate("BTC", "UAH")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, rate)
}

func (service *MainService) PostEmail(c *gin.Context) {
	request := c.Request
	writer := c.Writer
	headerContentType := request.Header.Get("Content-Type")

	if headerContentType != "application/x-www-form-urlencoded" {
		writer.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	err := request.ParseForm()

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	newEmail := request.FormValue("email")
	errSave := service.storageRepository.Save(newEmail)

	if errSave.Err != nil {
		writer.WriteHeader(service.storageErrorMapper.MapError(errSave))
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (service *MainService) GetEmails(c *gin.Context) {
	emails, err := service.storageRepository.AllEmails()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, emails)
}

func (service *MainService) SendEmails(c *gin.Context) {
	emails, err := service.storageRepository.AllEmails()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rate, err := service.externalService.CurrentRate("BTC", "UAH")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = service.notificationService.Send(emails, rate)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
