package client

import (
	"net/http"

	"go.bytebuilders.dev/client/api"
)

func (c *Client) GetCurrentUser() (*api.User, error) {
	var user api.User
	apiPath := "/user"
	err := c.getParsedResponse(http.MethodGet, apiPath, jsonHeader, nil, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
