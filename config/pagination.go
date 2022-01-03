package config

import "goblog/pkg/config"

func init()  {
	config.Add("pagination", config.StrMap{
		// 默认每页条数
		"perpage": 10,

		// URL 中用以分表多少页的参数
		"url_query": "page",
	})
}
