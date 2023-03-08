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
