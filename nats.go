package client

import (
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
)

func (c *Client) NewNatsConnection(name string) (*nats.Conn, error) {
	apiPath := "/user/nats/credentials"

	res := struct {
		NatsEndpoints []string `json:"natsEndpoints"`
		Credentials   []byte   `json:"credentials"`
	}{}
	err := c.getParsedResponse(http.MethodGet, apiPath, jsonHeader, nil, &res)
	if err != nil {
		return nil, err
	}

	if name == "" {
		name = "b3-client-go"
	}
	credFile, err := os.CreateTemp("", "nats-*.creds")
	if err != nil {
		return nil, err
	}

	_, err = credFile.Write(res.Credentials)
	if err != nil {
		return nil, err
	}
	defer os.Remove(credFile.Name())

	return nats.Connect(
		res.NatsEndpoints[0],
		nats.Name(name),
		nats.UserCredentials(credFile.Name()),
		nats.NoReconnect(),
	)
}
