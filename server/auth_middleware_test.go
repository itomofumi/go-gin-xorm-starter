package server_test

import (
	"testing"

	"github.com/gemcook/go-gin-xorm-starter/server"
)

func TestGetBearer(t *testing.T) {
	type args struct {
		auth []string
	}
	tests := []struct {
		name    string
		args    args
		wantJwt string
		wantOk  bool
	}{
		// TODO: Add test cases.
		{"ok", args{[]string{"Bearer eyJraW"}}, "eyJraW", true},
		{"ng", args{[]string{""}}, "", false},
		{"wrong format", args{[]string{"Bearer Bearer eyJraW"}}, "", false},
		{"space before token", args{[]string{" Bearer zzz "}}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJwt, gotOk := server.GetBearer(tt.args.auth)
			if gotJwt != tt.wantJwt {
				t.Errorf("GetBearer() gotJwt = %v, want %v", gotJwt, tt.wantJwt)
			}
			if gotOk != tt.wantOk {
				t.Errorf("GetBearer() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
