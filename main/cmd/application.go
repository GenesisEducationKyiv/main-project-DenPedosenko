package cmd

import "github.com/gin-gonic/gin"

type Application struct {
	Router Router[*gin.Engine]
}

func NewApplication() *Application {
	return &Application{
		Router: NewGinRouter(NewInitializer().initialize()),
	}
}
