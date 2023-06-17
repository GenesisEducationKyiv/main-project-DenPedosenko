package main

import (
	"bufio"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type User interface {
	GetEmail() string
}

type user struct {
	email string
}

func (u *user) GetEmail() string {
	return u.email
}

func NewUser(email string) User {
	return &user{email: email}
}

func getRate(c *gin.Context) {
	user := NewUser("test")
	_ = user.GetEmail()
	rate, err := getCurrentBTCToUAHRate()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, rate)
}

func postEmail(c *gin.Context) {
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
	httpStatus, errSave := saveEmailToStorage(newEmail)
	if errSave != nil {
		writter.WriteHeader(httpStatus)
		return
	}

	writter.WriteHeader(httpStatus)
}

func getEmails(c *gin.Context) {
	file, err := openFile(os.O_RDONLY)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	var emails []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		emails = append(emails, scanner.Text())
	}
	c.IndentedJSON(http.StatusOK, emails)
}

func sendEmails(c *gin.Context) {
	file, err := openFile(os.O_RDONLY)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	var emails []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		emails = append(emails, scanner.Text())
	}

	err = send(emails)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
