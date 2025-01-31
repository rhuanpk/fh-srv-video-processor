package router

import (
	"net/http"
	"strings"
)

func hasWildcard(route string) bool {
	return strings.Contains(route, "*")
}

func matchRoute(route, requested string) bool {
	prefix := strings.TrimSuffix(route, "*")
	return strings.HasPrefix(requested, prefix)
}

func wildcardHandler(route string, handler http.HandlerFunc) http.HandlerFunc {
	if !hasWildcard(route) {
		return handler
	}

	return func(res http.ResponseWriter, req *http.Request) {
		if !matchRoute(route, req.URL.Path) {
			http.NotFound(res, req)
			return
		}
		handler(res, req)
	}
}
