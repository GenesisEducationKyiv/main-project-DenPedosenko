package config

import (
	"context"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type configLoader struct {
	filepath string
}

func NewConfigLoader(filepath string) *configLoader {
	return &configLoader{filepath: filepath}
}

func (loader *configLoader) loadConfig() (*config, error) {
	content, err := os.ReadFile(loader.filepath)
	if err != nil {
		return nil, err
	}

	var conf config
	err = yaml.Unmarshal(content, &conf)

	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func (loader *configLoader) GetContext() (context.Context, error) {
	config, err := loader.loadConfig()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return WithConfig(context.Background(), config), nil
}
