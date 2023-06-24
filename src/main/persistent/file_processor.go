package persistent

import (
	"os"
)

const defaultFilePermission = 0o644

type FileProcessor interface {
	OpenFile(ac int) (*os.File, error)
	CreateFile() (*os.File, error)
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
	file, err := os.OpenFile(fp.defaultFilePath, os.O_APPEND|ac, defaultFilePermission)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = fp.CreateFile()
			if err != nil {
				return nil, err
			}

			return file, nil
		}

		return nil, err
	}

	return file, nil
}

func (fp *FileProcessorImpl) CreateFile() (*os.File, error) {
	file, err := os.Create(fp.defaultFilePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
