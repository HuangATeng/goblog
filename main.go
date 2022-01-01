package main

import (
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	//"database/sql"
	"goblog/config"
	c "goblog/pkg/config"
	"net/http"
)

//var router *mux.Router
//var db *sql.DB






func init()  {
	// 初始化配置信息
	config.Initialize()
}


func main() {

	// 初始化数据库
	//database.Initialize()
	//db = database.DB

	//route.Initialize()
	//router = route.Router
	bootstrap.SetupDB()
	router := bootstrap.SetupRoute()

	// url 后缀处理
	http.ListenAndServe(":" + c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
	//logger.LogError(err)
	//http.ListenAndServe(":3000", removeTrailingSlash(router))
}