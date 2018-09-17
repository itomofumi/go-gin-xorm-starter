package repository_test

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/gemcook/go-gin-xorm-starter/infra"
	"github.com/go-xorm/xorm"
)

var dockerMySQLContainerName = "go-gin-xorm-starter-test-mysql"
var dockerMySQLPort = "11336"

// Setup initializes test environment.
// Call cleanup func with 'defer'.
func Setup(t *testing.T) (engine *xorm.Engine, cleanup func()) {

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Chdir("..")
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("docker", "container", "run",
		"--rm",
		"-p", dockerMySQLPort+":3306",
		"-e", "MYSQL_ROOT_PASSWORD=password",
		"-e", "TZ=Asia/Tokyo",
		"mysql:5.7.21")

	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	setMySQLTestEnv()
	initDatabase(t)

	engine, err = infra.InitMySQLEngine(infra.LoadMySQLConfigEnv())
	if err != nil {
		t.Fatal(err)
	}

	engine.SetConnMaxLifetime(time.Second)

	// clean up function.
	return engine, func() {
		defer os.Chdir(currentDir)
		engine.Close()

		// send interrupt signal to docker command.
		cmd.Process.Signal(syscall.SIGINT)
	}
}

func setMySQLTestEnv() {
	os.Setenv("DATABASE_HOST", "localhost:"+dockerMySQLPort)
	os.Setenv("DATABASE_NAME", "go_gin_xorm_starter")
	os.Setenv("DATABASE_USER", "root")
	os.Setenv("DATABASE_PASSWORD", "password")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_DIR", "log/test")
}

func initDatabase(t *testing.T) {
	mysqlConf := infra.LoadMySQLConfigEnv()
	mysqlConf.DBName = ""
	connStr := mysqlConf.FormatDSN()
	fmt.Println(connStr)
	err := infra.RunSQLFile(connStr, "./fixtures/db.sql")
	if err != nil {
		t.Fatal(err)
	}
}
