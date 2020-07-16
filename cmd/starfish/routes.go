package main

import (
	"net/http"

	"github.com/bjjb/starfish"
)

type route struct {
	host    string
	handler http.Handler
}

func build(rule *rule) *route {
	switch rule.Action {
	case "s", "serve":
		return &route{rule.Host, starfish.Static(rule.Target)}
	case "f", "forward":
		return &route{rule.Host, starfish.Proxy(rule.Target)}
	default:
		return &route{rule.Host, starfish.BadGateway}
	}
}

func (route *route) Match(r *http.Request) bool {
	return r.Host == route.host
}

func (route *route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route.handler.ServeHTTP(w, r)
}
