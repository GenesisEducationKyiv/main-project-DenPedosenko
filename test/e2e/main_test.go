package e2e

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"testing"
)

type header struct {
	header string
	value  string
}

var contentTypeHeader = header{header: "Content-Type", value: "application/x-www-form-urlencoded"}

const (
	baseURL            = "http://localhost:8080/api"
	statusErrorMessage = "Response status should be %s"
	status200          = "200 OK"
	status409          = "409 Conflict"
	email              = "test@gmail.com"
)

func TestE2EMain(t *testing.T) {
	client := resty.New()

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				var s = initialize()

				router := gin.Default()
				router.GET("api/rate", s.GetRate)
				router.GET("api/subscribe", s.GetEmails)
				router.POST("api/subscribe", s.PostEmail)
				router.POST("api/sendEmails", s.SendEmails)

				err := router.Run("localhost:8080")

				if err != nil {
					log.Fatal(err)
				}

				func() {
					err := os.Remove("emails.txt")
					if err != nil {
						log.Fatal(err)
					}
				}()
			}
		}
	}(ctx)

	type e2eTestCase struct {
		name           string
		expectedStatus string
		method         string
		url            string
	}

	for _, scenario := range []e2eTestCase{
		{
			name:           "shouldGetRate",
			expectedStatus: status200,
			method:         resty.MethodGet,
			url:            fmt.Sprintf("%s/rate", baseURL),
		},
		{
			name:           "shouldGetAllEmails",
			expectedStatus: status200,
			method:         resty.MethodGet,
			url:            fmt.Sprintf("%s/subscribe", baseURL),
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			resp, err := client.R().Execute(scenario.method, scenario.url)

			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, scenario.expectedStatus, resp.Status(), fmt.Sprintf(statusErrorMessage, scenario.expectedStatus))
			assert.NotEmpty(t, resp.String())
		})
	}

	var testBody = fmt.Sprintf("email=%s", email)

	t.Run("shouldPostEmail", func(t *testing.T) {
		resp, err := client.R().
			SetHeader(contentTypeHeader.header, contentTypeHeader.value).
			SetBody(testBody).
			Post(fmt.Sprintf("%s/subscribe", baseURL))
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, status200, resp.Status(), fmt.Sprintf(statusErrorMessage, status200))
	})

	t.Run("shouldReturnConflict", func(t *testing.T) {
		resp, err := client.R().
			SetHeader(contentTypeHeader.header, contentTypeHeader.value).
			SetBody(testBody).
			Post(fmt.Sprintf("%s/subscribe", baseURL))
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, status409, resp.Status(), fmt.Sprintf(statusErrorMessage, status409))
	})

	t.Run("shouldSendEmails", func(t *testing.T) {
		resp, err := client.R().Post(fmt.Sprintf("%s/sendEmails", baseURL))
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, status200, resp.Status(), fmt.Sprintf(statusErrorMessage, status200))
	})
	cancel()
}
