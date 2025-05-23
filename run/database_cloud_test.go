package run

import (
	"testing"

	"github.com/beego/beego"
)

func initConfig() {
	err := beego.LoadAppConfig("ini", "../conf/app.conf")
	if err != nil {
		panic(err)
	}

	beego.BConfig.WebConfig.Session.SessionOn = true
}

func TestCreateDatabaseCloud(t *testing.T) {
	initConfig()
	InitRdsClient()

	dbName := "casibase_test3"
	_, err := gitCreateDatabaseCloud(dbName)
	if err != nil {
		panic(err)
	}

	err = addDatabaseUser(dbName)
	if err != nil {
		panic(err)
	}
}
