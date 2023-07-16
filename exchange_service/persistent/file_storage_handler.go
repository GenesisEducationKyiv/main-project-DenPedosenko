package persistent

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type FileStorage struct {
	fileProcessor FileProcessor
	logger        logger
}

type StorageError struct {
	Err  error
	Code ErrorCode
}

type ErrorCode int

const (
	Conflict     ErrorCode = 0
	UnknownError ErrorCode = 1
)

func NewFileStorage(fileProcessor FileProcessor, logger logger) *FileStorage {
	return &FileStorage{
		fileProcessor: fileProcessor,
		logger:        logger,
	}
}

func (storage *FileStorage) Save(email string) StorageError {
	file, err := storage.fileProcessor.OpenFile(os.O_WRONLY)

	if err != nil {
		storage.logger.Error(fmt.Sprintf("File storage handler error: %s", err.Error()))

		return StorageError{
			Err:  errors.New("can't open file"),
			Code: UnknownError,
		}
	}

	if storage.IsEmailAlreadyExists(email) {
		storage.logger.Info(fmt.Sprintf("Email already exists: %s", email))

		return StorageError{
			errors.New("email already exists"),
			Conflict,
		}
	}

	_, errWrite := file.WriteString(email + "\n")

	if errWrite != nil {
		storage.logger.Error(fmt.Sprintf("File storage handler error: %s", errWrite.Error()))

		return StorageError{
			errors.New("can't write to file"),
			UnknownError,
		}
	}

	defer file.Close()

	return StorageError{
		Err:  nil,
		Code: -1,
	}
}

func (storage *FileStorage) IsEmailAlreadyExists(newEmail string) bool {
	file, err := storage.fileProcessor.OpenFile(os.O_RDONLY)

	if err != nil {
		storage.logger.Error(fmt.Sprintf("File storage handler error: %s", err.Error()))
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
	file, err := storage.fileProcessor.OpenFile(os.O_RDONLY)

	if err != nil {
		storage.logger.Error(fmt.Sprintf("File storage handler error: %s", err.Error()))
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
