package router

import (
	"fmt"
	"log"
	"net/http"

	"extractor/pkg/router/middle"
)

func Setup() {
	log.Println("api base route:", BaseRoute)

	for _, scheme := range Schemes {
		if scheme.Method == "" {
			scheme.Method = methodAny
		}

		for index, parse := range scheme.Parses() {
			route := scheme.Routes[index]

			log.Printf(
				`%s%`+fmt.Sprint((4+(7-len(scheme.Method))))+`s%s`,
				scheme.Method, " -> ", route,
			)

			http.HandleFunc(parse,
				middle.Wrapper(
					wildcardHandler(route, scheme.Handler),
					scheme.Middlewares...,
				),
			)
		}
	}

	// fs := http.FileServer(http.Dir(consts.ImagesFolderRelative))
	// http.Handle("/image/", http.StripPrefix("/image", fs))
}
