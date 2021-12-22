package middlewares

import (
	"goblog/pkg/session"
	"net/http"
)

// StartSession 开启 session 会话控制
func StartSession(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 启动会话
		session.StartSession(w, r)

		// 继续处理接下请求
		next.ServeHTTP(w, r)
	})
}