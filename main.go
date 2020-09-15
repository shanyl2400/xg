package main

import (
	"os"
	"xg/conf"
	"xg/da"
	"xg/log"
	_ "xg/log"
	"xg/route"
)

func loadConfig() {
	dbConn := os.Getenv("xg_db_conn")
	if dbConn == "" {
		panic("xg_db_conn env is required")
	}
	redisConn := os.Getenv("xg_redis_conn")
	if redisConn == "" {
		panic("xg_redis_conn env is required")
	}
	uploadPath := os.Getenv("xg_upload_path")
	if uploadPath == "" {
		panic("xg_upload_path env is required")
	}

	logPath := os.Getenv("xg_log_path")

	c := &conf.Config{
		DBConnectionString: dbConn,
		RedisConnectionString: redisConn,
		LogPath: logPath,
		UploadPath: uploadPath,
	}
	conf.Set(c)
}

// @title xg REST API
// @version 0.1 alpha
// @description  xg backend rest api
// @termsOfService https://localhost:8088/v1
func main() {
	engine := route.Get()
	loadConfig()

	//初始化日志
	log.LogInit()

	//迁移数据库
	da.AutoMigrate()
	da.InitData(true)

	engine.Run(":8088")
}
