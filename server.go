package main

import (
	"fmt"

	"github.com/yaling888/soulx/common"
)

func reloadConfig() error {
	c, err := common.Parse()
	if err != nil {
		return err
	}

	req := common.MakeRequest(c)

	force := "true"

	body := map[string]interface{}{
		"payload": "",
		"path":    "",
	}

	fail := common.HTTPError{}

	resp, err := req.R().SetError(&fail).SetBody(&body).SetQueryParam("force", force).Put("/configs")

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("request failed: %s", fail.Message)
	}

	return nil
}
