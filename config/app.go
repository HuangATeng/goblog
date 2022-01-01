package config

import "goblog/pkg/config"

func init()  {
	config.Add("app", config.StrMap{
		// 应用名称，暂时未使用
		"name": config.Env("APP_NAME", "GoBlog"),
		// 当前环境 用以区分多环境
		"env": config.Env("APP_ENV", "production"),
		// 是否进入 debug 模式
		"debug": config.Env("APP_DEBUG", false),
		 // 应用服务端口
		 "port": config.Env("APP_PORT", "3000"),

		 // gorilla/session 在 Cookie中加密数据使用
		 "key": config.Env("APP_KEY","33446a9dcf9ea060a0a6532b166da32f304af0df"),
	})
}