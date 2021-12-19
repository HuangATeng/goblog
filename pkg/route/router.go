package route

import (
	"github.com/gorilla/mux"
	"net/http"
)

//// Router 路由对象
var route *mux.Router

func SetRoute(r *mux.Router)  {
	route = r
}
//
//// Initialize 初始化路由
//func Initialize()  {
//	Router = mux.NewRouter()
//	routes.RegisterWebRoutes(Router)
//}

// RouteName2URL 通过路由名称获取 URL
func Name2Url(routeName string, pairs ...string)  string{
	url, err := route.Get(routeName).URL(pairs...)
	if err != nil {
		return ""
	}

	return url.String()
}

// GetRouteVariable 获取 URI 路由参数
func GetRouteVariable(parameters string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameters]
}
