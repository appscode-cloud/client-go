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
	"reflect"
	"testing"

	"go.bytebuilders.dev/client-go"
	"go.bytebuilders.dev/client-go/api"
)

func TestClient_VerifyLicense(t *testing.T) {
	type fields struct {
		url         string
		accessToken string
		license     string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *api.License
		wantErr bool
	}{
		{
			name: "InvalidLicenseVerification",
			fields: fields{
				url:     "http://localhost:3000",
				license: "itsa.jwt.token",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ValidLicenseVerification",
			fields: fields{
				url:     "http://localhost:3000/",
				license: "eyJhbGci...OiJSUzIIn0.eyJhdWQiOlsib25sZS1...c3ViIjoiOCJ9.ZsDshA739P...OQzh6UOynsA",
			},
			want: &api.License{
				Status: "active",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewClient(tt.fields.accessToken, tt.fields.license, tt.fields.url)
			got, err := c.VerifyLicense()
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyLicense() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if got == nil || got.Status != tt.want.Status {
					t.Errorf("VerifyLicense() got = %v, want %v", got, tt.want)
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyLicense() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetLicensePlan(t *testing.T) {
	type fields struct {
		url         string
		accessToken string
		license     string
	}
	type args struct {
		clusterID      string
		productID      string
		productOwnerID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "InvalidLicense",
			fields: fields{
				url:         "http://localhost:3000",
				accessToken: "",
				license:     "itsa.jwt.token",
			},
			args: args{
				clusterID:      "not-a-id",
				productID:      "not-a-product-id",
				productOwnerID: 0,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "ValidLicense",
			fields: fields{
				url:         "http://localhost:3000",
				accessToken: "",
				license:     "eyJhbGci...OiJSUzIIn0.eyJhdWQiOlsib25sZS1...c3ViIjoiOCJ9.ZsDshA739P...OQzh6UOynsA",
			},
			args: args{
				clusterID:      "only-name-and-number",
				productID:      "prod_valid_id",
				productOwnerID: 8,
			},
			want:    "plan_GnwGFLUOMYlmfJ",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewClient(tt.fields.accessToken, tt.fields.license, tt.fields.url)
			got, err := c.GetLicensePlan(tt.args.clusterID, tt.args.productID, tt.args.productOwnerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLicensePlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLicensePlan() got = %v, want %v", got, tt.want)
			}
		})
	}
}
