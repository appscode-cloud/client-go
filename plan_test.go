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

	"go.bytebuilders.dev/client"
	"x-helm.dev/apimachinery/apis/products/v1alpha1"
)

func TestClient_GetProductPlans(t *testing.T) {
	type fields struct {
		url         string
		accessToken string
	}
	type args struct {
		productID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    v1alpha1.PlanList
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewClient(tt.fields.url).WithAccessToken(tt.fields.accessToken)
			_, err := c.GetProductPlans(tt.args.productID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProductPlans() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_GetProductPlan(t *testing.T) {
	type fields struct {
		url         string
		accessToken string
	}
	type args struct {
		productID string
		planID    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1alpha1.Plan
		wantErr bool
	}{
		{
			name: "Get product",
			fields: fields{
				url: "http://localhost:3000",
			},
			args: args{
				productID: "prod_valid_id",
				planID:    "plan_valid_id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewClient(tt.fields.url).WithAccessToken(tt.fields.accessToken)
			got, err := c.GetProductPlan(tt.args.productID, tt.args.planID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProductPlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("GetProductPlan got = nil")
			}
		})
	}
}
