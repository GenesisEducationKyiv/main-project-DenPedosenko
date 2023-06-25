package service

import (
	"net/http"

	"ses.genesis.com/exchange-web-service/src/main/notification"
	"ses.genesis.com/exchange-web-service/src/main/persistent"

	"github.com/gin-gonic/gin"
)

type MainService struct {
	externalService     ExternalService
	persistentService   persistent.PersistentStorage
	notificationService notification.NotifyService
}

func NewMainService(externalService ExternalService, persistentService persistent.PersistentStorage,
	notificationService notification.NotifyService) *MainService {
	return &MainService{
		externalService:     externalService,
		persistentService:   persistentService,
		notificationService: notificationService}
}

func (service *MainService) GetRate(c *gin.Context) {
	rate, err := service.externalService.GetCurrentBTCToUAHRate()
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
	httpStatus, errSave := service.persistentService.SaveEmailToStorage(newEmail)

	if errSave != nil {
		writer.WriteHeader(httpStatus)
		return
	}

	writer.WriteHeader(httpStatus)
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
	rate, _ := service.externalService.GetCurrentBTCToUAHRate()

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
