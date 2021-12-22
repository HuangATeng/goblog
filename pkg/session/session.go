package session

import (
	"github.com/gorilla/sessions"
	"goblog/pkg/logger"
	"net/http"
)

// Store gorilla sessions 的存储库

var Store = sessions.NewCookieStore([]byte("33446a9dcf9ea060a0a6532b166da32f304af0de"))

// Session 当前会话
var Session *sessions.Session

// Request 用以获取会话
var Request *http.Request

// Response 用以写入会话
var Response http.ResponseWriter

// StartSession 初始化会话， 在中间件中调用
func StartSession(w http.ResponseWriter, r *http.Request)  {
	var err error

	// gorilla/sessions 支持多会话 ，本项目使用单一会话
	Session, err = Store.Get(r, "goblog-session")
	logger.LogError(err)

	Request 	= r
	Response 	= w
}


// Put 写入键值对应会话数据
func Put(key string, value interface{})  {
	Session.Values[key] = value

	Save()
}

// Get 获取会话数据
func Get(key string) interface{}  {
	return Session.Values[key]

}

// Forget 删除会话
func Forget(key string)  {
	delete(Session.Values, key)
	Save()
}

// Flush 删除当前会话
func Flush()  {
	Session.Options.MaxAge = -1
	Save()
}

// Save 保持会话
func Save()  {
	// 非 HTTPS 的链接无法使用Secure 和 HttpOnly ,浏览器会报错
	// Session.Options.Secure = true
	// Session.Options.HttpOnly = true

	err := Session.Save(Request, Response)
	logger.LogError(err)
}