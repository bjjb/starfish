package starfish

import (
	"net/http"
)

// A Handler should have a Handle function which returns a Route
type Handler interface {
	Handle(http.Handler) Route
	HandleFunc(func(http.ResponseWriter, *http.Request)) Route
}

// A HandlerFunc is an adapter to allow the use of ordinary functions as
// Handlers.
type HandlerFunc func(http.Handler) Route

// Handle calls f(h).
func (f HandlerFunc) Handle(h http.Handler) Route {
	return f(h)
}

// HandleFunc calls f(http.HandlerFunc(h))
func (f HandlerFunc) HandleFunc(h func(http.ResponseWriter, *http.Request)) Route {
	return f(http.HandlerFunc(h))
}
