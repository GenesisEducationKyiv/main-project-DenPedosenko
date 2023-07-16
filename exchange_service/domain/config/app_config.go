package config

import (
	"context"
)

type AppConfig struct {
	EmailUser     string       `yaml:"username"`
	EmailPassword string       `yaml:"password"`
	EmailHost     string       `yaml:"host"`
	EmailPort     string       `yaml:"port"`
	LoggerConfig  LoggerConfig `yaml:"logger"`
	CoinAPI       ConfigAPI    `yaml:"coin_api"`
	CoinGecko     ConfigAPI    `yaml:"coin_gecko"`
	KuCoin        ConfigAPI    `yaml:"ku_coin"`
}

type ConfigAPI struct {
	Key string `yaml:"key"`
	URL string `yaml:"url"`
}

type LoggerConfig struct {
	URL string `yaml:"url"`
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
