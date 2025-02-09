package handler

import (
	"net/http"

	"extractor/internal/resouce/video"
)

type controller struct {
	service video.Service
}

func NewController(service video.Service) video.Controller {
	return &controller{service: service}
}

func (c controller) Process(res http.ResponseWriter, req *http.Request) {
	// zipsPaths, err := c.service.Process([]string{config.VideoPath}, config.FrameInterval, true)
	// zipsPaths, err := c.service.Process(config.VideosPaths, config.FrameInterval, true)
	// if err != nil {
	// 	log.Println("error in process videos:", err)
	// 	return
	// }

	// log.Println("success in process videos:", zipsPaths)
}
