package common

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type HTTPError struct {
	Message string `json:"message"`
}

func MakeRequest(c *Config) *resty.Client {
	client := resty.New().SetBaseURL("http://" + c.ExternalController)

	if c.Secret != "" {
		client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.Secret))
	}

	return client
}
