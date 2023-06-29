package persistent

import (
	"os"
)

const defaultFilePermission = 0o644

type FileProcessor interface {
	OpenFile(ac int) (*os.File, error)
}

type FileProcessorImpl struct {
	defaultFilePath string
}

func NewFileProcessor(path string) *FileProcessorImpl {
	return &FileProcessorImpl{
		defaultFilePath: path,
	}
}

func (fp *FileProcessorImpl) OpenFile(ac int) (*os.File, error) {
	file, err := os.OpenFile(fp.defaultFilePath, os.O_APPEND|os.O_CREATE|ac, defaultFilePermission)
	if err != nil {
		return nil, err
	}

	return file, nil
}
