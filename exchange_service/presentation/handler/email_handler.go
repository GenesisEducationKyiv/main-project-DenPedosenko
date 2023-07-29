package handler

import (
	"exchange-web-service/persistent"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StorageErrorMapper[T any, R any] interface {
	MapError(code T) R
}

type logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Close()
}

type EmailHandler struct {
	storageRepository  StorageRepository
	storageErrorMapper StorageErrorMapper[persistent.StorageError, int]
	logger             logger
}

type StorageRepository interface {
	AllEmails() ([]string, error)
	Save(email string) persistent.StorageError
	Remove(email string) persistent.StorageError
	IsEmailAlreadyExists(newEmail string) bool
}

func NewEmailHandler(storageRepository StorageRepository,
	storageErrorMapper StorageErrorMapper[persistent.StorageError, int], logger logger) *EmailHandler {
	return &EmailHandler{
		storageRepository:  storageRepository,
		storageErrorMapper: storageErrorMapper,
		logger:             logger,
	}
}

func (h *EmailHandler) PostEmail(c *gin.Context) {
	request := c.Request
	writer := c.Writer
	headerContentType := request.Header.Get("Content-Type")

	if headerContentType != "application/x-www-form-urlencoded" {
		h.logger.Error("Unsupported media type")
		writer.WriteHeader(http.StatusUnsupportedMediaType)

		return
	}

	err := request.ParseForm()

	if err != nil {
		h.logger.Info(fmt.Sprintf("Bad request: %s", request.Body))
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	newEmail := request.FormValue("email")
	errSave := h.storageRepository.Save(newEmail)

	if errSave.Err != nil {
		h.logger.Error(fmt.Sprintf("Error saving email: %s", errSave.Err.Error()))
		writer.WriteHeader(h.storageErrorMapper.MapError(errSave))

		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (h *EmailHandler) GetEmails(c *gin.Context) {
	emails, err := h.storageRepository.AllEmails()
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error getting emails: %s", err.Error()))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.IndentedJSON(http.StatusOK, emails)
}
