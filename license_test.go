package bytebuilders

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestVerifyLicense(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    *License
		wantErr bool
	}{
		//{
		//	name: "Invalid license",
		//	args: args{
		//		token: "fkdsja.afdskjafldjajka.afdkjakdjagjk.hdsfjkag",
		//	},
		//	want:    nil,
		//	wantErr: true,
		//},
		{
			name: "Valid license",
			args: args{
				token: "eyJhbGciOiJSUzI1NiIsImtpZCI6IlRETU5PTUlCIiwidHlwIjoiSldUIn0.eyJhdWQiOlsiMWVmMzNiNTMtZjRmZi00MjU2LWE0ZmYtNzgzOGJlZTA3YjMxIl0sImV4cCI6MTU3NTQ2MTE3NSwiaWF0IjoxNTczMDQxNDM1LCJpc3MiOiJieXRlLmJ1aWxkZXJzIiwianRpIjoiYTZiNDFlNGItOGRlNC00YjdjLWIyNTgtNzQ5NWJjZDFmNmZlIiwibmJmIjoxNTcyODY5MTc1LCJzaWkiOiJzaV9HN0l2Q1lSNVZMWERKRiIsInN1YiI6IjgifQ.kODX62cMpcjdNlJotuUSXCHUuZHfqgpTkydrG7kaI_2sl62niFOdvV8ILVUIsHdMqU6hMV-2N8JUyY20Na5sSOG6HKMi194H290XusvfJ8S0JLZUe8IFaBfRMZkegL0ZLWxywHO9x1UPYv2GGPeJic-wkFU0HLFrFY9tdIeQhaQ12glVFzoGqmv4Mu2WuI4ZxcI_9LnlOlggFb5B-M6092DTFBTk_So4JJsDinV6XG2Z-M3RtPTz0y-ezUzXf3xf6uyaNCUyKxt6KhTH9uJpu4qK7ZApxmr-n8ABKsVhTQkxyRoL_odcO7gj-JA79Cuf8Kv8y988FcFP_b8LlTG3lw",
			},
			want:    &License{},
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

			if got != nil {
				spew.Dump(got.SubscribedPlans)
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("VerifyLicense() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
