/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
