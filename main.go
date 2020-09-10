package main

import (
	"os"
	"xg/conf"
	"xg/da"
	_ "xg/log"
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

// @title xg REST API
// @version 0.1 alpha
// @description  xg backend rest api
// @termsOfService https://localhost:8088/v1
func main() {
	engine := route.Get()
	loadConfig()

	//迁移数据库
	da.AutoMigrate()
	da.InitData(true)

	engine.Run(":8088")
}
