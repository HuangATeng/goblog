package config

import "goblog/pkg/config"

func init()  {
	config.Add("session", config.StrMap{
		// 目前仅支持 Cookie
		"default": config.Env("SESSION_DRIVER", "cookie"),

		// 会话 cookie 名称
		"session_name": config.Env("SESSION_NAME", "goblog-session"),
	})
}
