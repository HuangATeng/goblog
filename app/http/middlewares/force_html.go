package middlewares

import "net/http"

// html标头中间件 强制标头返回 HTML 内容类型
func ForceHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}
