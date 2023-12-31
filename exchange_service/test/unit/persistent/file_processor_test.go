package persistent_test

import (
	"exchange-web-service/persistent"
	testservice "exchange-web-service/test/unit/service"
	"os"
	"testing"
)

const TestFilePath = "test_file_storage.txt"

var testFileProcessor = persistent.NewFileProcessor(TestFilePath, testservice.TestLogger{})

func TestOpenFile(t *testing.T) {
	t.Run("should openFileIfExist", func(t *testing.T) {
		_, _ = os.Create(TestFilePath)
		file, err := testFileProcessor.OpenFile(os.O_RDWR)
		if err != nil {
			t.Error("Expected error to be nil")
		}

		if file == nil {
			t.Error("Expected file to be not nil")
		}
		_ = os.Remove(TestFilePath)
	})

	t.Run("should createFileIfNotExist", func(t *testing.T) {
		file, err := testFileProcessor.OpenFile(os.O_RDWR)
		if err != nil {
			t.Error("Expected error to be nil")
		}

		if file == nil {
			t.Error("Expected file to be not nil")
		}

		err = os.Remove(TestFilePath)

		if err != nil {
			t.Error("Expected error to be nil")
		}
	})

	t.Run("shouldThrowErrorIfCanCreateFile", func(t *testing.T) {
		testFileProcessor = persistent.NewFileProcessor("", testservice.TestLogger{})
		file, err := testFileProcessor.OpenFile(os.O_RDWR)
		if err == nil {
			t.Error("Expected error to be not nil")
		}

		if file != nil {
			t.Error("Expected file to be nil")
		}
	})
}
