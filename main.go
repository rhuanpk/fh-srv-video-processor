package main

import (
	"log"
	"net/http"
	"time"

	"extractor/internal/config"
	_ "extractor/internal/infra/api"
)

func main() {
	defer func(start time.Time) {
		log.Println("execution time:", time.Since(start).String())
	}(time.Now())

	log.Printf("http server listnig on %q\n", config.APIPort.Parse())
	log.Fatalln(
		"error in start http server:",
		http.ListenAndServe(config.APIPort.Parse(), nil),
	)
}
