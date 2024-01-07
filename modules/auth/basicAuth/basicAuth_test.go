package basicauth

import (
	"fmt"
	"testing"
)

func Test_hashPassword(t *testing.T) {
	type args struct {
		password string
		hashType HashType
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "測試未設定加密",
			args: args{
				password: "ttttest",
				hashType: "",
			},
			want:    "ttttest",
			wantErr: false,
		},
		{
			name: "測試SHA512",
			args: args{
				password: "ttttest",
				hashType: HashTypeSHA512,
			},
			want:    "56f97b4795cfb61ad630746c9607c267ad0e20af78d0ad3abb1560910b343e1163a306508ba10b3f0af4be649dc504d9d7a9cffb80cdbf95b570904f3b6ce90b",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.UseEncodeType = tt.args.hashType
			got, err := HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("hashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("hashPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	type args struct {
		password      string
		truePsassword string
		hashType      HashType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "測試SHA512",
			args: args{
				password:      "talence",
				truePsassword: "ccaf2f2763f13b07600304525b1ed78e7cb96ef82fcb40e1f44fd27f9fdd2091e222c91e8ab2d28db1d5bc0aff447e3c3d19fba616ce748bf1e99183d7e8f914",
				hashType:      HashTypeSHA512,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.UseEncodeType = tt.args.hashType
			cfg.PasswordHash = "at104"
			got, _ := HashPassword(tt.args.password)
			fmt.Println(got)
			if err := Check(tt.args.password, tt.args.truePsassword); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
