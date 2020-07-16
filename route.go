package starfish

import "net/http"

// A Route is a http.Handler and a Matcher
type Route interface {
	Matcher
	http.Handler
}

// Simple route implementation
type route struct {
	matcher Matcher
	handler http.Handler
}

// Match with the route's matcher
func (route *route) Match(r *http.Request) bool {
	return route.matcher.Match(r)
}

// ServeHTTP with the route's handler
func (route *route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route.handler.ServeHTTP(w, r)
	return
}
