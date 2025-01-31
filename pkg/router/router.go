package router

import (
	"net/http"

	"extractor/pkg/router/middle"
)

type Scheme struct {
	Method      string
	Routes      []string
	Handler     http.HandlerFunc
	Middlewares []middle.Middleware
}

const (
	BaseRoute = "/api/v1"
	methodAny = "ANY"
)

var Schemes []Scheme

func (s Scheme) buildRoute(index int) string {
	if index < 0 || index >= len(s.Routes) {
		return ""
	}
	route := BaseRoute + s.Routes[index]
	if s.Method != methodAny {
		route = s.Method + " " + route
	}
	return route
}

func (s Scheme) Parse(index int) string {
	return s.buildRoute(index)
}

func (s Scheme) Parses() (parses []string) {
	for index := range s.Routes {
		parses = append(parses, s.buildRoute(index))
	}
	return
}

// AddScheme add a [Scheme] of router to set a new endpoint on API start. In
// [Scheme.Middlewares] field you must set the middlewares in ascendent order as
// you desire.
func AddScheme(scheme Scheme) {
	Schemes = append(Schemes, scheme)
}
