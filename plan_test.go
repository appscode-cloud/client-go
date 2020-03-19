package client_test

import (
	"net/http"
	"testing"

	"go.bytebuilders.dev/client-go"
	"kubepack.dev/kubepack/apis/kubepack/v1alpha1"
)

func TestClient_GetProductPlans(t *testing.T) {
	type fields struct {
		url         string
		accessToken string
		license     string
		sudo        string
		client      *http.Client
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
			c := client.NewClient(tt.fields.accessToken, tt.fields.license, tt.fields.url)
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
		license     string
		sudo        string
		client      *http.Client
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
			c := client.NewClient(tt.fields.accessToken, tt.fields.license, tt.fields.url)
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
