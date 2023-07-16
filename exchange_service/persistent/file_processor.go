package persistent

import (
	"fmt"
	"os"
)

const defaultFilePermission = 0o644

type FileProcessor interface {
	OpenFile(ac int) (*os.File, error)
}

type logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Close()
}

type FileProcessorImpl struct {
	defaultFilePath string
	logger          logger
}

func NewFileProcessor(path string, logger logger) *FileProcessorImpl {
	return &FileProcessorImpl{
		defaultFilePath: path,
		logger:          logger,
	}
}

func (fp *FileProcessorImpl) OpenFile(ac int) (*os.File, error) {
	file, err := os.OpenFile(fp.defaultFilePath, os.O_APPEND|os.O_CREATE|ac, defaultFilePermission)

	if err != nil {
		fp.logger.Error(err.Error())
		return nil, err
	}

	fp.logger.Debug(fmt.Sprintf("File opened from path: %s", fp.defaultFilePath))

	return file, nil
}
