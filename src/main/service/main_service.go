package service

import (
	"net/http"

	"ses.genesis.com/exchange-web-service/src/service/errormapper"

	"ses.genesis.com/exchange-web-service/src/notification"
	"ses.genesis.com/exchange-web-service/src/persistent"
	"ses.genesis.com/exchange-web-service/src/main/notification"
	"ses.genesis.com/exchange-web-service/src/main/persistent"

	"github.com/gin-gonic/gin"
)

type MainService struct {
	externalService     ExternalService
	persistentService   persistent.PersistentStorage
	notificationService notification.NotifyService
	persistentService   persistent.Storage
	notificationService notification.NotificationService
	storageErrorMapper  errormapper.StorageErrorMapper[persistent.StorageError, int]
}

func NewMainService(externalService ExternalService, persistentService persistent.Storage,
	notificationService notification.NotificationService,
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

	errParse := request.ParseForm()

	if errParse != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	newEmail := request.FormValue("email")
	errSave := service.persistentService.SaveEmailToStorage(newEmail)

	if errSave != nil {
		writer.WriteHeader(service.storageErrorMapper.MapError(*errSave))
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
