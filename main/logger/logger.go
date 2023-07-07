package logger

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Decorator struct {
}

func NewLogger() *Decorator {
	return &Decorator{}
}

func (logger *Decorator) NewLogResponseDecorator(client *resty.Client) *resty.Client {
	client.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		fmt.Printf("URL: %s |Status: %s | Response: %s\n", response.Request.URL, response.Status(), response.String())
		return nil
	})

	return client
}
