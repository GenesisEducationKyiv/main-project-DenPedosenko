package service

import (
	"net/http"

	"ses.genesis.com/exchange-web-service/src/notification"
	"ses.genesis.com/exchange-web-service/src/persistent"

	"github.com/gin-gonic/gin"
)

type mainService struct {
	externalService     ExternalService
	persistentService   persistent.PersistentStorage
	notificationService notification.NotificationService
}

func NewMainService(externalService ExternalService, persistentService persistent.PersistentStorage,
	notificationService notification.NotificationService) *mainService {
	return &mainService{
		externalService:     externalService,
		persistentService:   persistentService,
		notificationService: notificationService}
}

func (service *mainService) GetRate(c *gin.Context) {
	rate, err := service.externalService.GetCurrentBTCToUAHRate()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, rate)
}

func (service *mainService) PostEmail(c *gin.Context) {
	request := c.Request
	writter := c.Writer
	headerContentType := request.Header.Get("Content-Type")

	if headerContentType != "application/x-www-form-urlencoded" {
		writter.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	errParse := request.ParseForm()

	if errParse != nil {
		writter.WriteHeader(http.StatusBadRequest)
		return
	}

	newEmail := request.FormValue("email")
	httpStatus, errSave := service.persistentService.SaveEmailToStorage(newEmail)

	if errSave != nil {
		writter.WriteHeader(httpStatus)
		return
	}

	writter.WriteHeader(httpStatus)
}

func (service *mainService) GetEmails(c *gin.Context) {
	emails, err := service.persistentService.AllEmails()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, emails)
}

func (service *mainService) SendEmails(c *gin.Context) {
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
