package api

import (
	"net/http"

	processor "extractor/internal/resouce/processor/handler"
	zipper "extractor/internal/resouce/zipper/handler"
	"extractor/pkg/router"
	"extractor/pkg/router/middle"
)

func setupRun() {
	service := processor.NewService(zipper.NewService())
	controller := processor.NewController(service)
	router.AddScheme(router.Scheme{
		Method:      http.MethodGet,
		Routes:      []string{"/run"},
		Handler:     controller.Process,
		Middlewares: []middle.Middleware{middle.CORS},
	})
}
