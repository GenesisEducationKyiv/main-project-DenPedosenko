package cmd

import (
	"exchange-web-service/domain/config"
	"exchange-web-service/domain/rabittmq"

	"github.com/gin-gonic/gin"
)

const (
	configPath      = "./resources/application.yaml"
	fileStoragePath = "./resources/emails.txt"
)

type Application struct {
	Router Router[*gin.Engine]
	Logger *rabittmq.Logger
}

func NewApplication() *Application {
	configLoader := config.NewConfigLoader(configPath)
	ctx, err := configLoader.GetContext()

	logger, _ := rabittmq.NewLogger(config.GetConfigFromContext(ctx).LoggerConfig.URL)

	if err != nil {
		panic(err)
	}

	return &Application{
		Router: NewGinRouter(NewInitializer().initialize(ctx, logger)),
		Logger: logger,
	}
}
