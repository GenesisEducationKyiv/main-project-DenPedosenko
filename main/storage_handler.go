package main

import (
	"bufio"
	"errors"
	"net/http"
	"os"
)

const defaultFilePermission = 0o644
const defaultFilePath = "emails.txt"

func saveEmailToStorage(email string) (int, error) {
	file, err := openFile(os.O_WRONLY)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if isEmailAlreadyExists(email) {
		return http.StatusConflict, errors.New("email already exists")
	}

	_, errWrite := file.WriteString(email + "\n")

	if errWrite != nil {
		return http.StatusInternalServerError, err
	}

	defer file.Close()

	return http.StatusOK, nil
}

func openFile(ac int) (*os.File, error) {
	file, err := os.OpenFile(defaultFilePath, os.O_APPEND|ac, defaultFilePermission)
	if err != nil {
		if os.IsNotExist(err) {
			errFileCreation := createFile()
			if errFileCreation != nil {
				return nil, errFileCreation
			}

			return file, nil
		}

		return nil, err
	}

	return file, nil
}

func createFile() error {
	file, err := os.Create(defaultFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func isEmailAlreadyExists(newEmail string) bool {
	file, err := openFile(os.O_RDONLY)
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
