package main

import (
	"fmt"
	"os"

	"github.com/gemcook/go-gin-xorm-starter/infra"
	"github.com/gemcook/go-gin-xorm-starter/util"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

func runSource(mysqlConnectionString, sourceFilePath string) error {

	var err error
	engine, err := xorm.NewEngine("mysql", mysqlConnectionString)
	if err != nil {
		return err
	}

	engine.ShowSQL(true)
	engine.Logger().SetLevel(core.LOG_UNKNOWN)

	err = engine.Ping()
	if err != nil {
		return err
	}

	file, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = engine.Import(file)

	if err != nil {
		if err.Error() == "not an error" {
			err = nil
		} else {
			return err
		}
	}

	return nil
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("params error")
		return
	}

	util.LoadEnv()
	mysqlConf := infra.LoadMySQLConfigEnv()
	mysqlConf.DBName = ""
	mysqlConnectionString := mysqlConf.FormatDSN()
	fmt.Println(mysqlConnectionString)

	err := runSource(mysqlConnectionString, os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
}
