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
	defer os.Remove(credFile.Name()) // nolint:errcheck

	return nats.Connect(
		res.NatsEndpoints[0],
		nats.Name(name),
		nats.UserCredentials(credFile.Name()),
		nats.NoReconnect(),
	)
}
