package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SignInParams struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

func (c *Client) Signin(params SignInParams) ([]http.Cookie, error) {
	apiPath := "/user/signin"
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	_, cookies, err := c.getResponseWithCookies(http.MethodPost, apiPath, jsonHeader, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return cookies, nil
}

func (c *Client) Signout() error {
	apiPath := "/user/signout"
	_, err := c.getResponse(http.MethodGet, apiPath, jsonHeader, nil)
	if err != nil {
		return err
	}
	return nil
}
