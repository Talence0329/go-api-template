package auth

import (
	"backend/basic/apiprotocol"
	"fmt"
	"reflect"
	"testing"

	"github.com/lithammer/shortuuid/v4"
)

func Test_accountLogin(t *testing.T) {
	type args struct {
		account  string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    apiprotocol.Code
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, got, err := loginAccount(tt.args.account, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("accountLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("accountLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uuid(t *testing.T) {
	fmt.Println(shortuuid.New())
	t.Log(shortuuid.New())
}
