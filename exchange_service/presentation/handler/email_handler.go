package handler

import (
	"exchange-web-service/persistent"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StorageErrorMapper[T any, R any] interface {
	MapError(code T) R
}

type EmailHandler struct {
	storageRepository  StorageRepository
	storageErrorMapper StorageErrorMapper[persistent.StorageError, int]
}

type StorageRepository interface {
	AllEmails() ([]string, error)
	Save(email string) persistent.StorageError
	IsEmailAlreadyExists(newEmail string) bool
}

func NewEmailHandler(storageRepository StorageRepository,
	storageErrorMapper StorageErrorMapper[persistent.StorageError, int]) *EmailHandler {
	return &EmailHandler{
		storageRepository:  storageRepository,
		storageErrorMapper: storageErrorMapper,
	}
}

func (h *EmailHandler) PostEmail(c *gin.Context) {
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
	errSave := h.storageRepository.Save(newEmail)

	if errSave.Err != nil {
		writer.WriteHeader(h.storageErrorMapper.MapError(errSave))
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (h *EmailHandler) GetEmails(c *gin.Context) {
	emails, err := h.storageRepository.AllEmails()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, emails)
}
