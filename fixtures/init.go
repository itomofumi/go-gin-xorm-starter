package main

import (
	"fmt"
	"os"

	"github.com/gemcook/go-gin-xorm-starter/infra"
	"github.com/gemcook/go-gin-xorm-starter/util"
)

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

	err := infra.RunSQLFile(mysqlConnectionString, os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
}
