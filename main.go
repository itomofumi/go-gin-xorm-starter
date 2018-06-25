package main

import (
	"os"

	"github.com/gemcook/go-gin-xorm-starter/server"
	"github.com/gemcook/go-gin-xorm-starter/util"
)

func main() {
	util.LoadEnv()
	err := server.Start()
	if err != nil {
		util.GetLogger().Errorln(err)
		os.Exit(1)
	}
}
