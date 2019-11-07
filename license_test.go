package client

import (
	"reflect"
	"testing"

	"go.bytebuilder.dev/client/api"
)

func TestVerifyLicense(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    *api.License
		wantErr bool
	}{
		{
			name: "Invalid license",
			args: args{
				token: "fkdsja.afdskjafldjajka.afdkjakdjagjk",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Valid license",
			args: args{
				token: "eyJhbGciOiJS....idHlwIjoiSldUIn0.eyJhdWQiOlsiMWVmM.....RKRiIsInN1YiI6IjgifQ.kODX62cMpcjdNlJotuUSXC.....8FcFP_b8LlTG3lw",
			},
			want:    &api.License{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyLicense(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyLicense() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				if reflect.DeepEqual(got, nil) {
					t.Errorf("VerifyLicense() got = %v, don't want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGetLicensePlan(t *testing.T) {
	type args struct {
		token          string
		clusterID      string
		productID      string
		productOwnerID int64
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
	}{
		{
			name: "InvalidGetLicensePlan",
			args: args{
				token:          "jsaflkdjas.fkasjdlakj.afdskja;lsf",
				clusterID:      "kakjfksljdfkl-aflkdjak",
				productID:      "akfdsjaklj",
				productOwnerID: 0,
			},
			want:  false,
			want1: "",
		},
		{
			name: "ValidGetLicensePlan",
			args: args{
				token:          "eyJhbGciOiJS....idHlwIjoiSldUIn0.eyJhdWQiOlsiMWVmM.....RKRiIsInN1YiI6IjgifQ.kODX62cMpcjdNlJotuUSXC.....8FcFP_b8LlTG3lw",
				clusterID:      "1ef33b53-f4ff-4256-a4ff-7838bee07b31",
				productID:      "prod_F..........Lb9",
				productOwnerID: 8,
			},
			want:  true,
			want1: "plan_G3..........v5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetLicensePlan(tt.args.token, tt.args.clusterID, tt.args.productID, tt.args.productOwnerID)
			if got != tt.want {
				t.Errorf("GetLicensePlan() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetLicensePlan() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
