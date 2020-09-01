package main

import (
	"os"
	"xg/conf"
	"xg/da"
	"xg/route"
)

func loadConfig() {
	dbConn := os.Getenv("xg_db_conn")
	if dbConn == "" {
		panic("xg_db_conn env is required")
	}
	c := &conf.Config{DBConnectionString: dbConn}
	conf.Set(c)
}

func main() {
	engine := route.Get()
	loadConfig()

	//迁移数据库
	da.AutoMigrate()
	da.InitData(false)

	engine.Run(":8088")
}
