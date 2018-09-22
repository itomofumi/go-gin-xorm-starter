package repository_test

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/gemcook/go-gin-xorm-starter/infra"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var dockerMySQLImage = "mysql:5.7.21"
var dockerMySQLPort = "11336"
var dockerMySQLName = "go-gin-xorm-starter-mysql" + dockerMySQLPort

// Setup initializes test environment.
// Call cleanup func with 'defer'.
func Setup(t *testing.T) (engine *xorm.Engine, cleanup func()) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker command is not installed")
	}

	dockerInfoCmd := exec.Command("docker", "info")
	err := dockerInfoCmd.Run()
	if err != nil {
		t.Skipf("docker daemon is not running. error=%v", err)
	}

	removeMySQLDockerContainer()

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Chdir("..")
	if err != nil {
		t.Fatal(err)
	}

	setMySQLTestEnv()
	mysqlConf := infra.LoadMySQLConfigEnv()

	dockerRunCmd := exec.Command("docker", "container", "run",
		"--rm",
		"--name", dockerMySQLName,
		"-p", dockerMySQLPort+":3306",
		"-e", "MYSQL_ROOT_PASSWORD="+mysqlConf.Passwd,
		"-e", "TZ=Asia/Tokyo",
		dockerMySQLImage)

	err = dockerRunCmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	err = waitDbIsReady()
	if err != nil {
		t.Fatal(err)
	}

	initDatabase(t)

	engine, err = infra.InitMySQLEngine(mysqlConf)
	if err != nil {
		t.Fatal(err)
	}
	engine.ShowSQL(false)
	engine.SetConnMaxLifetime(time.Millisecond * 1)

	// clean up function.
	return engine, func() {
		if t.Skipped() {
			return
		}

		defer os.Chdir(currentDir)
		engine.Close()

		// send interrupt signal to docker command.
		err = dockerRunCmd.Process.Signal(syscall.SIGINT)
		if err != nil {
			fmt.Println("SIGINT:", err)
		}

		err = dockerRunCmd.Wait()
		if err != nil {
			fmt.Println("dockerRunCmd.Wait()", err)
		}

		removeMySQLDockerContainer()
	}
}

func setMySQLTestEnv() {
	os.Setenv("DATABASE_HOST", "0.0.0.0:"+dockerMySQLPort)
	os.Setenv("DATABASE_NAME", "go-gin-xorm-starter")
	os.Setenv("DATABASE_USER", "root")
	os.Setenv("DATABASE_PASSWORD", "password")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_DIR", "log/test")
}

func removeMySQLDockerContainer() {
	dockerRmImageCmd := exec.Command("sh", "-c", fmt.Sprintf(`'docker container rm -f $(docker container ps -q -f "name=%s")'`, dockerMySQLName))

	err := dockerRmImageCmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func waitDbIsReady() error {
	mysqlConf := infra.LoadMySQLConfigEnv()
	mysqlConf.DBName = ""
	engine, err := infra.InitMySQLEngine(mysqlConf)
	if err != nil {
		return err
	}

	engine.SetConnMaxLifetime(time.Millisecond * 1)
	engine.ShowSQL(false)
	engine.Logger().SetLevel(core.LOG_WARNING)

	retry := 10
	for i := 0; i < retry; i++ {
		err := engine.Ping()
		if err == nil {
			return nil
		}
		fmt.Println(err)
		time.Sleep(time.Duration(100 * time.Millisecond))
	}
	return fmt.Errorf("cannot connect to host: %v", mysqlConf.Addr)
}

func initDatabase(t *testing.T) {
	mysqlConf := infra.LoadMySQLConfigEnv()
	mysqlConf.DBName = ""
	connStr := mysqlConf.FormatDSN()
	err := infra.RunSQLFile(connStr, "./fixtures/db.sql")
	if err != nil {
		t.Fatal(err)
	}
}
