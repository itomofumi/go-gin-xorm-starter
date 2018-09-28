package repository_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/moby/moby/client"

	"github.com/gemcook/go-gin-xorm-starter/infra"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var dockerMySQLImage = "mysql:5.7.21"
var dockerMySQLPort = "23306"
var dockerMySQLName = "go-gin-xorm-starter-mysql" + dockerMySQLPort

func TestMain(m *testing.M) {
	if _, err := exec.LookPath("docker"); err != nil {
		fmt.Println("docker command is not installed")
		os.Exit(0)
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("failed to initialize docker client")
		os.Exit(0)
	}

	_, err = cli.Info(context.TODO())
	if err != nil {
		fmt.Println(fmt.Errorf("docker daemon is not running. error=%v", err))
		os.Exit(0)
	}

	cleanup := setupInfra()

	ret := m.Run()

	cleanup()

	os.Exit(ret)
}

func setupInfra() (cleanup func()) {
	setMySQLTestEnv()
	mysqlConf := infra.LoadMySQLConfigEnv()

	cli, err := client.NewEnvClient()
	if err != nil {
		panic("failed to initialize docker client")
	}

	// start mysql if not running.
	list := listMySQLDockerContainers(cli)
	if len(list) == 0 {
		pullMySQLDockerImage(cli)

		created, err := cli.ContainerCreate(context.TODO(), &container.Config{
			Image:        dockerMySQLImage,
			ExposedPorts: nat.PortSet{nat.Port("3306"): struct{}{}},
			Env: []string{
				"MYSQL_ROOT_PASSWORD=" + mysqlConf.Passwd,
				"TZ=Asia/Tokyo",
			},
			AttachStdout: true,
		}, &container.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port("3306"): []nat.PortBinding{{HostPort: dockerMySQLPort}},
			},
		}, &network.NetworkingConfig{}, dockerMySQLName)

		if err != nil {
			panic(fmt.Errorf("docker container create failed: %v", err))
		}

		err = cli.ContainerStart(context.TODO(), created.ID, types.ContainerStartOptions{})
		if err != nil {
			panic(fmt.Errorf("docker container start failed: %v", err))
		}

		reader, err := cli.ContainerLogs(context.TODO(), dockerMySQLName, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: false,
			Follow:     true,
			Tail:       "1000",
		})

		if err != nil {
			panic(fmt.Errorf("docker container logs failed: %v", err))
		}

		sc := bufio.NewScanner(reader)
		var timeoutSec int = 30
		timeoutCh := time.After(time.Second * time.Duration(timeoutSec))
	waitMySQL:
		for {
			select {
			case <-timeoutCh:
				panic(fmt.Errorf("mysql initialization timeout %v sec", timeoutSec))
			default:
				if sc.Scan() && strings.Contains(sc.Text(), "ready for connections") {
					fmt.Println("mysql: ready for connections")
					break waitMySQL
				}
			}
		}

		// wait additional few seconds.
		time.Sleep(time.Second * 3)

		// and check connection.
		err = waitDbIsReady()

		if err != nil {
			panic("cannot connect to mysql")
		}
	}

	return func() {
		// currently DO NOT remove mysql container for usability to iterate test.
		// removeMySQLDockerContainer(cli, listMySQLDockerContainers(cli))
	}
}

// Setup initializes mysql database for each case.
// Call cleanup func with 'defer'.
func setupDB(t *testing.T) (engine *xorm.Engine, cleanup func()) {
	t.Helper()

	err := waitDbIsReady()
	if err != nil {
		t.Fatalf("waitDbIsReady() failed. %v", err)
	}

	initDatabase(t)

	mysqlConf := infra.LoadMySQLConfigEnv()
	engine, err = infra.InitMySQLEngine(mysqlConf)
	if err != nil {
		t.Fatal(err)
	}
	engine.ShowSQL(false)
	engine.SetConnMaxLifetime(time.Second * 1)

	// clean up function.
	return engine, func() {
		engine.Close()
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

func pullMySQLDockerImage(cli *client.Client) {
	filterMap := map[string][]string{"reference": {dockerMySQLImage}}
	filterBytes, _ := json.Marshal(filterMap)
	filter, err := filters.FromParam(string(filterBytes))
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(context.TODO(), types.ImageListOptions{
		All:     true,
		Filters: filter,
	})
	if err != nil {
		panic(err)
	}

	if len(images) == 0 {
		reader, err := cli.ImagePull(context.TODO(), dockerMySQLImage, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
		sc := bufio.NewScanner(reader)
		for sc.Scan() {
			text := sc.Text()
			if strings.Contains(text, "Downloading") || strings.Contains(text, "Extracting") {
				continue
			}
			fmt.Println(text)
		}
	}
}

func listMySQLDockerContainers(cli *client.Client) []types.Container {
	filterMap := map[string][]string{"name": {dockerMySQLName}}
	filterBytes, _ := json.Marshal(filterMap)
	filter, _ := filters.FromParam(string(filterBytes))

	opts := types.ContainerListOptions{
		Filters: filter,
	}

	list, err := cli.ContainerList(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	return list
}

func removeMySQLDockerContainer(cli *client.Client, targets []types.Container) {
	for _, t := range targets {
		err := cli.ContainerRemove(context.TODO(), t.ID, types.ContainerRemoveOptions{
			Force: true,
		})

		if err != nil {
			panic(err)
		}
	}
}

func waitDbIsReady() error {
	mysqlConf := infra.LoadMySQLConfigEnv()
	mysqlConf.DBName = ""
	engine, err := infra.InitMySQLEngine(mysqlConf)
	if err != nil {
		return err
	}

	engine.SetConnMaxLifetime(time.Second * 1)
	engine.ShowSQL(false)
	engine.Logger().SetLevel(core.LOG_WARNING)

	retry := 15
	for i := 0; i < retry; i++ {
		err := engine.Ping()
		if err == nil {
			fmt.Println("db connection established")
			return nil
		}
		time.Sleep(time.Duration(1000 * time.Millisecond))
	}
	return fmt.Errorf("cannot connect to host: %v", mysqlConf.Addr)
}

func initDatabase(t *testing.T) {
	t.Helper()
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Chdir("..")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(currentDir)

	mysqlConf := infra.LoadMySQLConfigEnv()
	mysqlConf.DBName = ""
	connStr := mysqlConf.FormatDSN()
	err = infra.RunSQLFile(connStr, "./fixtures/db.sql")
	if err != nil {
		t.Fatal(err)
	}
}
