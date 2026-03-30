package utils

import "net/http"

type middleware func(http.Handler) http.HandlerFunc

func ChainMiddleWare(handle http.HandlerFunc, middlewares ...middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handle = middlewares[i](handle)
	}
	return handle
}
