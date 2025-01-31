package processor

import "net/http"

type Service interface {
	Process([]string, int, bool) ([]string, error)
}

type Controller interface {
	Process(http.ResponseWriter, *http.Request)
}
