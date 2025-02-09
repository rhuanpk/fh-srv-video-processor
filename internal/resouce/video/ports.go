package video

import "net/http"

type Service interface {
	Process(videosPaths []string, frameInterval int, highQuality bool) ([]string, error)
}

type Controller interface {
	Process(http.ResponseWriter, *http.Request)
}
