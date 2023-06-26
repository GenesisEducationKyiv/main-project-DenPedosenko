package persistent

import (
	"os"
)

const defaultFilePermission = 0o644
const defaultFilePath = "emails.txt"

type FileProcessor interface {
	openFile(ac int) (*os.File, error)
	createFile() error
}

type fileProcessor struct {
}

func NewFileProcessor() *fileProcessor {
	return &fileProcessor{}
}

func (fp *fileProcessor) openFile(ac int) (*os.File, error) {
	file, err := os.OpenFile(defaultFilePath, os.O_APPEND|os.O_CREATE|ac, defaultFilePermission)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (fp *fileProcessor) createFile() error {
	file, err := os.Create(defaultFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}
