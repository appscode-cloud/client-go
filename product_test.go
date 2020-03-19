/*
Copyright 2019 AppsCode Inc.

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

	"go.bytebuilders.dev/client-go"
)

func TestClient_GetProductByID(t *testing.T) {
	type fields struct {
		url         string
		accessToken string
		license     string
	}
	type args struct {
		productID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Get product",
			fields: fields{
				url: "http://localhost:3000",
			},
			args: args{
				productID: "prod_valid_id",
			},
			wantErr: false,
		},
		{
			name: "Do not get product",
			fields: fields{
				url: "http://localhost:3000",
			},
			args: args{
				productID: "invalid-id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewClient(tt.fields.accessToken, tt.fields.license, tt.fields.url)
			got, err := c.GetProductByID(tt.args.productID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProductByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("GetProductByID() got = nil, want")
			}
		})
	}
}

func TestClient_GetProductByOwnerAndKey(t *testing.T) {
	type fields struct {
		url         string
		accessToken string
		license     string
	}
	type args struct {
		owner string
		key   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Get product",
			fields: fields{
				url: "http://localhost:3000",
			},
			args: args{
				owner: "system-admin",
				key:   "kubedb",
			},
			wantErr: false,
		},
		{
			name: "Get product",
			fields: fields{
				url: "http://localhost:3000",
			},
			args: args{
				owner: "system-admin",
				key:   "not-kubedb",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewClient(tt.fields.accessToken, tt.fields.license, tt.fields.url)
			got, err := c.GetProductByOwnerAndKey(tt.args.owner, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProductByOwnerAndKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("GetProductByID() got = nil")
			}
		})
	}
}
