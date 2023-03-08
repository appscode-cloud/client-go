package client_test

import (
	"testing"

	"go.bytebuilders.dev/client"
)

func TestClient_GetCurrentUser(t *testing.T) {
	t.Run("Get current user", func(t *testing.T) {
		c := client.NewClient("http://api.bb.test:3003").WithBasicAuth("appscode", "password")
		user, err := c.GetCurrentUser()
		if err != nil {
			t.Error(err)
			return
		}
		if user.UserName != "appscode" {
			t.Errorf("Expected: appscode Got: %s", user.UserName)
		}
	})
}
