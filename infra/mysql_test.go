package infra_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/itomofumi/go-gin-xorm-starter/infra"
	"github.com/stretchr/testify/assert"
	"xorm.io/core"
)

func setupEnv() {
	os.Setenv("DATABASE_HOST", "1.2.3.4")
	os.Setenv("DATABASE_NAME", "foo")
	os.Setenv("DATABASE_USER", "bar")
	os.Setenv("DATABASE_PASSWORD", "fizzbuzz")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_DIR", "log")
}

func TestLoadMySQLConfigEnv(t *testing.T) {
	setupEnv()
	mysqlConf := infra.LoadMySQLConfigEnv()

	assert := assert.New(t)
	assert.Equal("1.2.3.4", mysqlConf.Addr)
	assert.Equal("foo", mysqlConf.DBName)
	assert.Equal("bar", mysqlConf.User)
	assert.Equal("fizzbuzz", mysqlConf.Passwd)

}

func TestInitMySQLEngine(t *testing.T) {
	setupEnv()
	mysqlConf := infra.LoadMySQLConfigEnv()

	assert := assert.New(t)

	engine, err := infra.InitMySQLEngine(mysqlConf)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal("mysql", engine.DriverName())
}

func TestEscapeMySQLString(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"drop table", `'; DROP TABLE users;`, `\'; DROP TABLE users;`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := infra.EscapeMySQLString(tt.value); got != tt.want {
				t.Errorf("EscapeMySQLString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseLogLevel(t *testing.T) {

	tests := []struct {
		value   string
		want    core.LogLevel
		wantErr bool
	}{
		{"fatal", core.LOG_ERR, false},
		{"error", core.LOG_ERR, false},
		{"panic", core.LOG_ERR, false},
		{"warning", core.LOG_WARNING, false},
		{"warn", core.LOG_WARNING, false},
		{"info", core.LOG_INFO, false},
		{"debug", core.LOG_DEBUG, false},
		{"other", core.LOG_DEBUG, true},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got, err := infra.ParseLogLevel(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLogLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
