package middle

import "net/http"

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("access-control-allow-origin", "*")
		res.Header().Set("access-control-allow-methods", "*")
		res.Header().Set("access-control-allow-headers", "*")
		if req.Method == http.MethodOptions {
			res.WriteHeader(http.StatusContinue)
		}
		next(res, req)
	}
}
