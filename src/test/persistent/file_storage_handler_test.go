package persistent_test

import (
	"errors"
	"os"
	"testing"

	"ses.genesis.com/exchange-web-service/src/main/persistent"
)

type TestFileProcessor struct {
}

type FailTestFileProcessor struct {
}

func (tfp *FailTestFileProcessor) OpenFile(_ int) (*os.File, error) {
	return nil, errors.New("test error")
}

func (tfp *FailTestFileProcessor) CreateFile() (*os.File, error) {
	return nil, errors.New("test error")
}

func (tfp *TestFileProcessor) OpenFile(_ int) (*os.File, error) {
	file, _ := os.OpenFile("test_file_storage.txt", os.O_APPEND|os.O_RDWR, 0666)
	return file, nil
}

func (tfp *TestFileProcessor) CreateFile() (*os.File, error) {
	return nil, nil
}

func TestFileStorage_AllEmails(t *testing.T) {
	beforeEach()
	t.Run("should return email from file", func(t *testing.T) {
		var fs = persistent.NewFileStorage(&TestFileProcessor{})
		emails, err := fs.AllEmails()

		if err != nil {
			t.Error("Expected error to be nil")
		}

		if emails[0] != "test@gmail.com" {
			t.Error("Expected email to be test@gmail.com")
		}
	})

	t.Run("should return nil if something goes wrong", func(t *testing.T) {
		beforeEach()
		var fs = persistent.NewFileStorage(&FailTestFileProcessor{})
		_, err := fs.AllEmails()

		if err == nil {
			t.Error("Expected error to be nil")
		}
	})
}

func TestFileStorage_SaveEmailToStorage(t *testing.T) {
	beforeEach()
	t.Run("should return OK if email is saved", func(t *testing.T) {
		var fs = persistent.NewFileStorage(&TestFileProcessor{})
		status, err := fs.SaveEmailToStorage("new_test_email")
		if err != nil {
			t.Error("Expected error to be nil")
		}

		if status != 200 {
			t.Error("Expected status to be 200")
		}
	})

	t.Run("should return 409 if email already exists", func(t *testing.T) {
		var fs = persistent.NewFileStorage(&TestFileProcessor{})
		status, err := fs.SaveEmailToStorage("test@gmail.com")
		if err == nil {
			t.Error("Expected error to be not nil")
		}

		if status != 409 {
			t.Error("Expected status to be 409")
		}
	})

	t.Run("should return 500 if something goes wrong", func(t *testing.T) {
		var fs = persistent.NewFileStorage(&FailTestFileProcessor{})
		status, err := fs.SaveEmailToStorage("test@gmail.com")

		if err == nil {
			t.Error("Expected error to be not nil")
		}

		if status != 500 {
			t.Error("Expected status to be 500")
		}
	})
}

func beforeEach() {
	_ = os.Remove("test_file_storage.txt")
	file, _ := os.OpenFile("test_file_storage.txt", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	file.WriteString("test@gmail.com\n") //nolint:errcheck
	file.Close()
}
