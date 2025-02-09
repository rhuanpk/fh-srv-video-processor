package api

import (
	"net/http"

	video "extractor/internal/resouce/video/handler"
	zipper "extractor/internal/resouce/zipper/handler"
	"extractor/pkg/router"
	"extractor/pkg/router/middle"
)

func setupRun() {
	service := video.NewService(zipper.NewService())
	controller := video.NewController(service)
	router.AddScheme(router.Scheme{
		Method:      http.MethodGet,
		Routes:      []string{"/run"},
		Handler:     controller.Process,
		Middlewares: []middle.Middleware{middle.CORS},
	})
}
