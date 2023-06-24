package config

import (
	"context"
)

type config struct {
	APIKey        string `yaml:"key"`
	APIURL        string `yaml:"url"`
	EmailUser     string `yaml:"username"`
	EmailPassword string `yaml:"password"`
	EmailHost     string `yaml:"host"`
	EmailPort     string `yaml:"port"`
}

type configKey struct{}

func WithConfig(ctx context.Context, config *config) context.Context {
	return context.WithValue(ctx, configKey{}, config)
}

func GetConfig(ctx context.Context) *config {
	if config, ok := ctx.Value(configKey{}).(*config); ok {
		return config
	}

	return nil
}
