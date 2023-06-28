package config

import (
	"context"
)

type AppConfig struct {
	APIKey        string `yaml:"key"`
	APIURL        string `yaml:"url"`
	EmailUser     string `yaml:"username"`
	EmailPassword string `yaml:"password"`
	EmailHost     string `yaml:"host"`
	EmailPort     string `yaml:"port"`
}

type configKey struct{}

func WithConfig(ctx context.Context, config *AppConfig) context.Context {
	return context.WithValue(ctx, configKey{}, config)
}

func GetConfigFromContext(ctx context.Context) *AppConfig {
	if config, ok := ctx.Value(configKey{}).(*AppConfig); ok {
		return config
	}

	return nil
}
