package middle

import "net/http"

func NoCache(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("cache-control", "no-cache, no-store, must-revalidate")
		next(res, req)
	}
}
