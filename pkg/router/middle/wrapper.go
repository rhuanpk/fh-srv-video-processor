package middle

import "net/http"

func Wrapper(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	handler = func(next http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			next(res, req)
		}
	}(handler)

	for index := len(middlewares) - 1; index > -1; index-- {
		handler = middlewares[index](handler)
	}
	return handler
}
