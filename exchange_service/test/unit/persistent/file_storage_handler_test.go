package persistent_test

import (
	"errors"
	"exchange-web-service/persistent"
	testservice "exchange-web-service/test/unit/service"
	"os"
	"testing"
)

type TestFileProcessor struct {
}

type FailTestFileProcessor struct {
}

func (tfp *FailTestFileProcessor) OpenFile(_ int) (*os.File, error) {
	return nil, errors.New("test error")
}

func (tfp *TestFileProcessor) OpenFile(_ int) (*os.File, error) {
	file, _ := os.OpenFile(TestFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	return file, nil
}

func (tfp *TestFileProcessor) RewriteFile(_ []byte) error {
	return nil
}

func (tfp *FailTestFileProcessor) RewriteFile(_ []byte) error {
	return nil
}

func TestFileStorage_AllEmails(t *testing.T) {
	cleanUpTestData()
	t.Run("should return email from file", func(t *testing.T) {
		var fs = persistent.NewFileStorage(&TestFileProcessor{}, testservice.TestLogger{})
		emails, err := fs.AllEmails()

		if err != nil {
			t.Error("Expected error to be nil")
		}

		if emails[0] != "test@gmail.com" {
			t.Error("Expected email to be test@gmail.com")
		}
	})

	t.Run("should return nil if something goes wrong", func(t *testing.T) {
		var fs = persistent.NewFileStorage(&FailTestFileProcessor{}, testservice.TestLogger{})
		_, err := fs.AllEmails()

		if err == nil {
			t.Error("Expected error to be nil")
		}

		defer cleanUpTestData()
	})
}

func TestFileStorage_SaveEmailToStorage(t *testing.T) {
	cleanUpTestData()
	t.Run("should return OK if email is saved", func(t *testing.T) {
		var fs = persistent.NewFileStorage(&TestFileProcessor{}, testservice.TestLogger{})
		err := fs.Save("new_test_email")
		if err.Err != nil {
			t.Error("Expected error to be nil")
		}
	})

	t.Run("should return error with code 0 if email already exists", func(t *testing.T) {
		var fs = persistent.NewFileStorage(&TestFileProcessor{}, testservice.TestLogger{})
		err := fs.Save("test@gmail.com")
		if err.Err == nil {
			t.Error("Expected error to be not nil")
		}

		if err.Code != 0 {
			t.Error("Expected code to be 0")
		}
	})

	t.Run("should return error with code 1 if something goes wrong", func(t *testing.T) {
		var fs = persistent.NewFileStorage(&FailTestFileProcessor{}, testservice.TestLogger{})
		err := fs.Save("test@gmail.com")

		if err.Err == nil {
			t.Error("Expected error to be not nil")
		}

		if err.Code != 1 {
			t.Error("Expected status to be 1")
		}
		defer cleanUpTestData()
	})
}

func cleanUpTestData() {
	_ = os.Remove(TestFilePath)
	file, err := os.OpenFile(TestFilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		panic(err)
	}

	file.WriteString("test@gmail.com\n") //nolint:errcheck
	file.Close()
}
