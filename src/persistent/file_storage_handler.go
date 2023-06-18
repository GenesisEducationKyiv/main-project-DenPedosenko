package persistent

import (
	"bufio"
	"errors"
	"net/http"
	"os"
)

type FileStorage struct {
	fileProcessor IFileProcessor
}

func NewFileStorage(fileProcessor IFileProcessor) PersistentStorage {
	return &FileStorage{
		fileProcessor: fileProcessor,
	}
}

func (storage *FileStorage) SaveEmailToStorage(email string) (int, error) {
	file, err := storage.fileProcessor.openFile(os.O_WRONLY)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if storage.IsEmailAlreadyExists(email) {
		return http.StatusConflict, errors.New("email already exists")
	}

	_, errWrite := file.WriteString(email + "\n")

	if errWrite != nil {
		return http.StatusInternalServerError, err
	}

	defer file.Close()

	return http.StatusOK, nil
}

func (storage *FileStorage) IsEmailAlreadyExists(newEmail string) bool {
	file, err := storage.fileProcessor.openFile(os.O_RDONLY)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == newEmail {
			return true
		}
	}

	return false
}

func (storage *FileStorage) AllEmails() ([]string, error) {
	file, err := storage.fileProcessor.openFile(os.O_RDONLY)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var emails []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		emails = append(emails, scanner.Text())
	}

	return emails, nil
}
