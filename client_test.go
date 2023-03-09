/*
Copyright 2020 AppsCode Inc.

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
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestServerURL     = "http://api.bb.test:3003"
	TestServerUser    = "appscode"
	TestServrPassword = "password"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Version",
			want: "v0.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Version(); got != tt.want {
				t.Errorf("Version() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		accessToken string
		baseURL     []string
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "NewClient",
			args: args{
				accessToken: "<a-valid-access-token",
			},
			want: &Client{
				url:         "https://byte.builders",
				accessToken: "<a-valid-access-token",
				client:      &http.Client{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.baseURL...).WithAccessToken(tt.args.accessToken); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClientWithHTTP(t *testing.T) {
	type args struct {
		httpClient  *http.Client
		accessToken string
		baseURL     []string
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "NewClientWithHTTP",
			args: args{
				httpClient:  nil,
				accessToken: "",
			},
			want: &Client{
				url: "https://byte.builders",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClientWithHTTP(tt.args.httpClient, tt.args.baseURL...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientWithHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SetHTTPClient(t *testing.T) {
	type fields struct {
		accessToken string
	}
	type args struct {
		client *http.Client
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "SettHTTPClient",
			fields: fields{},
			args:   args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient().WithAccessToken(tt.fields.accessToken)
			c.SetHTTPClient(tt.args.client)
		})
	}
}

func TestClient_getStatusCode(t *testing.T) {
	type fields struct {
		accessToken string
	}
	type args struct {
		method string
		path   string
		header http.Header
		body   io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "StatusCode",
			fields:  fields{},
			args:    args{},
			want:    http.StatusNotFound,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient().WithAccessToken(tt.fields.accessToken)
			got, err := c.getStatusCode(tt.args.method, tt.args.path, tt.args.header, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStatusCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getStatusCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_getParsedResponse(t *testing.T) {
	type fields struct {
		accessToken string
	}
	type args struct {
		method string
		path   string
		header http.Header
		body   io.Reader
		obj    interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "getParsedResponse",
			fields:  fields{},
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient().WithAccessToken(tt.fields.accessToken)
			if err := c.getParsedResponse(tt.args.method, tt.args.path, tt.args.header, tt.args.body, tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("getParsedResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_WithCookie(t *testing.T) {
	t.Run("Should be able to get user using cookie", func(t *testing.T) {
		c := NewClient(TestServerURL)
		cookies, err := c.Signin(SignInParams{
			UserName: TestServerUser,
			Password: TestServrPassword,
		})
		if !assert.Nil(t, err) {
			return
		}

		c = NewClient(TestServerURL).WithCookies(cookies)
		user, err := c.GetCurrentUser()
		if !assert.Nil(t, err) {
			return
		}
		if !assert.NotNil(t, user) {
			return
		}
		assert.Equal(t, user.UserName, TestServerUser)
	})
}
