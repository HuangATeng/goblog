package route

import "github.com/gorilla/mux"

// Router 路由对象
var Router *mux.Router

// Initialize 初始化路由
func Initialize()  {
	Router = mux.NewRouter()
}

// RouteName2URL 通过路由名称获取 URL
func Name2Url(routeName string, pairs ...string)  string{
	url, err := Router.Get(routeName).URL(pairs...)
	if err != nil {
		return ""
	}

	return url.String()
}
