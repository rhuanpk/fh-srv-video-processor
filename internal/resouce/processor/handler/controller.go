package handler

import (
	"log"
	"net/http"

	"extractor/internal/config"
	"extractor/internal/resouce/processor"
)

type controller struct {
	service processor.Service
}

func NewController(service processor.Service) processor.Controller {
	return &controller{service: service}
}

func (c controller) Process(res http.ResponseWriter, req *http.Request) {
	// zipsPaths, err := c.service.Process([]string{config.VideoPath}, config.FrameInterval, true)
	zipsPaths, err := c.service.Process(config.VideosPaths, config.FrameInterval, true)
	if err != nil {
		log.Println("error in process videos:", err)
		return
	}

	log.Println("success in process videos:", zipsPaths)
}
