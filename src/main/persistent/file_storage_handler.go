package persistent

import (
	"bufio"
	"errors"
	"os"
)

type FileStorage struct {
	fileProcessor FileProcessor
}

type StorageError struct {
	Error error
	Code  ErrorCode
}

type ErrorCode int

const (
	Conflict    ErrorCode = 0
	UnkownError ErrorCode = 1
)

func NewFileStorage(fileProcessor FileProcessor) Storage {
	return &FileStorage{
		fileProcessor: fileProcessor,
	}
}

func (storage *FileStorage) SaveEmailToStorage(email string) *StorageError {
	file, err := storage.fileProcessor.openFile(os.O_WRONLY)

	if err != nil {
		return &StorageError{
			errors.New("email already exists"),
			UnkownError,
		}
	}

	if storage.IsEmailAlreadyExists(email) {
		return &StorageError{
			errors.New("email already exists"),
			Conflict,
		}
	}

	file, err := storage.fileProcessor.OpenFile(os.O_WRONLY)

	if err != nil {
		return int(UnknownError), err
	}

	_, errWrite := file.WriteString(email + "\n")

	if errWrite != nil {
		return &StorageError{
			errors.New("email already exists"),
			UnkownError,
		}
	}

	defer file.Close()

	return nil
}

func (storage *FileStorage) IsEmailAlreadyExists(newEmail string) bool {
	file, err := storage.fileProcessor.OpenFile(os.O_RDONLY)
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
	file, err := storage.fileProcessor.OpenFile(os.O_RDONLY)
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
