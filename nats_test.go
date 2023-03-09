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

package client_test

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.bytebuilders.dev/client"
)

func TestNatsConnection(t *testing.T) {
	t.Run("Should connect with NATS server", func(t *testing.T) {
		c := client.NewClient(client.TestServerURL).WithBasicAuth(client.TestServerUser, client.TestServrPassword)
		nc, err := c.NewNatsConnection("unit-test")
		if !assert.Nil(t, err) {
			return
		}
		defer nc.Close()

		status := nc.Status()
		assert.Equal(t, status, nats.CONNECTED)
	})
}
