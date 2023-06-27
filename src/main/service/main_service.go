package service

import (
	"net/http"

	"ses.genesis.com/exchange-web-service/src/main/notification"
	"ses.genesis.com/exchange-web-service/src/main/persistent"
	"ses.genesis.com/exchange-web-service/src/main/service/errormapper"

	"github.com/gin-gonic/gin"
)

type MainService struct {
	externalService     ExternalService
	notificationService notification.NotifyService
	persistentService   persistent.Storage
	storageErrorMapper  errormapper.StorageErrorMapper[persistent.StorageError, int]
}

func NewMainService(externalService ExternalService, persistentService persistent.Storage,
	notificationService notification.NotifyService,
	storageErrorMapper errormapper.StorageErrorMapper[persistent.StorageError, int]) *MainService {
	return &MainService{
		externalService:     externalService,
		persistentService:   persistentService,
		notificationService: notificationService,
		storageErrorMapper:  storageErrorMapper}
}

func (service *MainService) GetRate(c *gin.Context) {
	rate, err := service.externalService.CurrentBTCToUAHRate()
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
	errSave := service.persistentService.SaveEmailToStorage(newEmail)

	if errSave.Error != nil {
		writer.WriteHeader(service.storageErrorMapper.MapError(errSave))
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (service *MainService) GetEmails(c *gin.Context) {
	emails, err := service.persistentService.AllEmails()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, emails)
}

func (service *MainService) SendEmails(c *gin.Context) {
	emails, err := service.persistentService.AllEmails()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rate, err := service.externalService.CurrentBTCToUAHRate()

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
