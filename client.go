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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var jsonHeader = http.Header{"Content-Type": []string{"application/json;charset=UTF-8"}}

// Version return the library version
func Version() string {
	return "v0.0.1"
}

// Client represents a ByteBuilders api client
type Client struct {
	url         string
	accessToken string
	client      *http.Client
	org         string
	username    string
	password    string
	cookies     []http.Cookie
}

// NewClient initializes and returns a API client.
func NewClient(baseURL ...string) *Client {
	url := "https://byte.builders"
	if len(baseURL) > 0 {
		url = baseURL[0]
	}
	return &Client{
		url:    strings.TrimSuffix(url, "/"),
		client: &http.Client{},
	}
}

func (c *Client) WithAccessToken(accessToken string) *Client {
	c.accessToken = accessToken
	return c
}

func (c *Client) WithBasicAuth(username, password string) *Client {
	c.username = username
	c.password = password
	return c
}

func (c *Client) WithCookies(cookies []http.Cookie) *Client {
	c.cookies = cookies
	return c
}

func (c *Client) WithOrganization(org string) *Client {
	c.org = org
	return c
}

// NewClientWithHTTP creates an API client with a custom http client
func NewClientWithHTTP(httpClient *http.Client, baseURL ...string) *Client {
	client := NewClient(baseURL...)
	client.client = httpClient
	return client
}

// SetHTTPClient replaces default http.Client with user given one.
func (c *Client) SetHTTPClient(client *http.Client) {
	c.client = client
}

func (c *Client) doRequest(method, path string, header http.Header, body io.Reader) (*http.Response, error) {
	path = strings.TrimPrefix(path, "/")
	req, err := http.NewRequest(method, c.url+"/api/v1/"+path, body)
	if err != nil {
		return nil, err
	}
	if len(c.accessToken) != 0 {
		req.Header.Set("Authorization", "token "+c.accessToken)
	}
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}
	if c.cookies != nil {
		c.addCookies(req)
	}
	for k, v := range header {
		req.Header[k] = v
	}

	return c.client.Do(req)
}

func (c *Client) addCookies(req *http.Request) {
	var csrfToken string
	for i := range c.cookies {
		if c.cookies[i].Name == "_csrf" {
			csrfToken = c.cookies[i].Value
		}
		req.AddCookie(&c.cookies[i])
	}
	req.Header.Set("X-Csrf-Token", csrfToken)
}

var (
	ErrUnAuthorized   = errors.New("401 Unauthorized")
	ErrForbidden      = errors.New("403 Forbidden")
	ErrNotFound       = errors.New("404 Not Found")
	ErrStatusConflict = errors.New("409 Conflict")
)

func (c *Client) getResponse(method, path string, header http.Header, body io.Reader) ([]byte, error) {
	resp, err := c.doRequest(method, path, header, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := parseStatusCode(resp.StatusCode, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) getParsedResponse(method, path string, header http.Header, body io.Reader, obj interface{}) error {
	data, err := c.getResponse(method, path, header, body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}

func (c *Client) getResponseWithCookies(method, path string, header http.Header, body io.Reader) ([]byte, []http.Cookie, error) {
	resp, err := c.doRequest(method, path, header, body)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if err := parseStatusCode(resp.StatusCode, data); err != nil {
		return nil, nil, err
	}
	var cookies []http.Cookie
	for _, c := range resp.Cookies() {
		cookie := *c
		cookies = append(cookies, cookie)
	}
	return data, cookies, nil
}

func (c *Client) getStatusCode(method, path string, header http.Header, body io.Reader) (int, error) {
	resp, err := c.doRequest(method, path, header, body)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func parseStatusCode(statusCode int, data []byte) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return ErrUnAuthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusConflict:
		return ErrStatusConflict
	case http.StatusUnprocessableEntity:
		return fmt.Errorf("422 Unprocessable Entity: %s", string(data))
	}

	if statusCode/100 != 2 {
		errMap := make(map[string]interface{})
		if err := json.Unmarshal(data, &errMap); err != nil {
			// when the JSON can't be parsed, data was probably empty or a plain string,
			// so we try to return a helpful error anyway
			return fmt.Errorf("Unknown API Error: %d %s", statusCode, string(data))
		}
		return errors.New(errMap["message"].(string))
	}

	return nil
}

func (c *Client) getOrganization() (string, error) {
	org := c.org
	if org == "" {
		user, err := c.GetCurrentUser()
		if err != nil {
			return "", fmt.Errorf("organization name wasn't provided. Automatic user detection failed with error: %w", err)
		}
		org = user.UserName
	}
	return org, nil
}

type queryParams struct {
	key   string
	value string
}

func setQueryParams(apiPath string, params []queryParams) (string, error) {
	u, err := url.Parse(apiPath)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for i := range params {
		q.Set(params[i].key, params[i].value)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
