package infra

import (
	"fmt"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

func TestKVSClient_GetStruct(t *testing.T) {
	type testObj struct {
		Name string
	}

	type fields struct {
		Conn redis.Conn
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"[success] get",
			fields{Conn: func() redis.Conn {
				c := redigomock.NewConn()
				c.Command("GET", "key1").Expect(`{"Name":"ok"}`)
				return c
			}()},
			args{"key1"},
			false,
		},
		{"[fail] get",
			fields{Conn: func() redis.Conn {
				c := redigomock.NewConn()
				c.Command("GET", "key1").ExpectError(fmt.Errorf("key does not exist"))
				return c
			}()},
			args{"key2"},
			true,
		},
		{"[fail] get, not connected",
			fields{Conn: nil},
			args{"key2"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc := &KVSClient{
				Conn: tt.fields.Conn,
			}
			obj := testObj{}
			if err := kc.GetStruct(tt.args.key, &obj); (err != nil) != tt.wantErr {
				t.Fatalf("KVSClient.GetStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && obj.Name != "ok" {
				t.Errorf("got testObj.Name = %q, want = %q", obj.Name, "ok")
			}
		})
	}
}

func TestKVSClient_SetStruct(t *testing.T) {
	type testObj struct {
		Name string
	}
	type fields struct {
		Conn redis.Conn
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"[success] set",
			fields{Conn: func() redis.Conn {
				c := redigomock.NewConn()
				c.Command("MULTI").Expect("ok")
				c.Command("SET").Expect("ok")
				c.Command("EXPIRE").Expect("ok")
				c.Command("EXEC").Expect("ok")
				return c
			}()},
			args{"key1", &testObj{Name: "value1"}},
			false,
		},
		{"[fail] set",
			fields{Conn: func() redis.Conn {
				c := redigomock.NewConn()
				c.Command("MULTI").Expect("ok")
				c.Command("SET").ExpectError(fmt.Errorf("failed to set"))
				c.Command("EXPIRE").Expect("ok")
				c.Command("EXEC").Expect("ok")
				return c
			}()},
			args{"key2", nil},
			true,
		},
		{"[fail] set, not connected",
			fields{Conn: nil},
			args{"key2", nil},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			kc := &KVSClient{
				Conn: tt.fields.Conn,
			}

			if err := kc.SetStruct(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Fatalf("KVSClient.SetStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
