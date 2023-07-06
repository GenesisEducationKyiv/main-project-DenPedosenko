package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ses.genesis.com/exchange-web-service/main/persistent"
	"ses.genesis.com/exchange-web-service/main/service/errormapper"
)

type EmailController struct {
	storageRepository  StorageRepository
	storageErrorMapper errormapper.StorageErrorMapper[persistent.StorageError, int]
}

type StorageRepository interface {
	AllEmails() ([]string, error)
	Save(email string) persistent.StorageError
	IsEmailAlreadyExists(newEmail string) bool
}

func NewEmailController(storageRepository StorageRepository,
	storageErrorMapper errormapper.StorageErrorMapper[persistent.StorageError, int]) *EmailController {
	return &EmailController{
		storageRepository:  storageRepository,
		storageErrorMapper: storageErrorMapper,
	}
}

func (controller *EmailController) PostEmail(c *gin.Context) {
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
	errSave := controller.storageRepository.Save(newEmail)

	if errSave.Err != nil {
		writer.WriteHeader(controller.storageErrorMapper.MapError(errSave))
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (controller *EmailController) GetEmails(c *gin.Context) {
	emails, err := controller.storageRepository.AllEmails()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, emails)
}
