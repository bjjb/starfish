package starfish

import (
	"net/http"
)

// A Router is essentially a safely accessible list of Routes that also
// implements http.Handler.
type Router []Route

// ServeHTTP uses the first matching route to respond to the http request.
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for route := range router.routes() {
		if route.Match(r) {
			route.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

// Match returns a Handler for the given matcher
func (router *Router) Match(matcher Matcher) Handler {
	return HandlerFunc(func(h http.Handler) Route {
		route := &route{matcher, h}
		router.Push(route)
		return route
	})
}

// MatchFunc returns a Handler for the matcher
func (router *Router) MatchFunc(matchFunc func(*http.Request) bool) Handler {
	return router.Match(MatcherFunc(matchFunc))
}

// Routes safely returns a copy of the Router's routes
func (router *Router) Routes() []Route {
	routes := []Route{}
	for route := range router.routes() {
		routes = append(routes, route)
	}
	return routes
}

// Append appends Routes to the Router
func (router *Router) Append(routes ...Route) {
	router.Replace(append(router.Routes(), routes...))
}

// Push pushes a route to the router
func (router *Router) Push(route Route) {
	router.Append(route)
}

// PopN removes n routes from the Router
func (router *Router) PopN(n int) []Route {
	routes := router.Routes()
	if len(routes) < n {
		return []Route{}
	}
	index := len(routes) - n
	remaining, popped := routes[:index], routes[index:]
	router.Replace(remaining)
	return popped
}

// Pop removes a Route from the Router (and returns it)
func (router *Router) Pop() Route {
	popped := router.PopN(1)
	if len(popped) == 0 {
		return nil
	}
	return popped[0]
}

// Replace replaces all the routes in the Router, and returns the old set
func (router *Router) Replace(newRoutes []Route) []Route {
	oldRoutes := router.Routes()
	*router = newRoutes
	return oldRoutes
}

// Clear removes all the routes in the router
func (router *Router) Clear() {
	router.Replace([]Route{})
}

// ListenAndServe listens on the given port and serves HTTP
func (router *Router) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, router)
}

// ListenAndServeTLS listens on the given port and serves HTTPS
func (router *Router) ListenAndServeTLS(addr, cert, key string) error {
	return http.ListenAndServeTLS(addr, cert, key, router)
}

func (router *Router) routes() chan Route {
	ch := make(chan Route)
	go func(routes []Route) {
		defer close(ch)
		if routes == nil {
			return
		}
		for _, route := range routes {
			ch <- route
		}
	}([]Route(*router))
	return ch
}

// ServeHTTP calls the default router's ServeHTTP
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defaultRouter.ServeHTTP(w, r)
}

// Match calls the default router's Match
func Match(matcher Matcher) Handler {
	return defaultRouter.Match(matcher)
}

// MatchFunc calls the default router's MatchFunc
var MatchFunc = defaultRouter.MatchFunc

// Routes of the default router
var Routes = defaultRouter.Routes

// Append routes to the default router
var Append = defaultRouter.Append

// Push a route to the default router
var Push = defaultRouter.Push

// PopN routes from the default router
var PopN = defaultRouter.PopN

// Pop a route from the default router
var Pop = defaultRouter.Pop

// Replace the routes of the default router
var Replace = defaultRouter.Replace

// Clear the routes of the default router
var Clear = defaultRouter.Clear

// ListenAndServe serves HTTP on the default router
var ListenAndServe = defaultRouter.ListenAndServe

// ListenAndServeTLS serves HTTPS on the default router
var ListenAndServeTLS = defaultRouter.ListenAndServeTLS

var defaultRouter *Router

func init() {
	defaultRouter = new(Router)
}
