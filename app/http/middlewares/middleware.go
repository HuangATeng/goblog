package middlewares

import "net/http"

// httpHandlerFunc 简写 func(http.ResponseWriter, *http.Request)
type HttpHandlerFunc func(http.ResponseWriter, *http.Request)
