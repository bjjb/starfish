package starfish

import "net/http"

// A Matcher can check a http.Request to see if it should serve it
type Matcher interface {
	Match(*http.Request) bool
}

// A MatcherFunc is an adapter to allow the use of ordinary functions as
// Matchers.
type MatcherFunc func(*http.Request) bool

// Match calls f(r).
func (f MatcherFunc) Match(r *http.Request) bool {
	return f(r)
}

// BoolMatcher gives a Matcher which always either matches or not, depending
// on its underlying value
type BoolMatcher bool

// Match implements Matcher for m
func (m BoolMatcher) Match(r *http.Request) bool {
	return bool(m)
}

// Always is a Matcher which always matches
const Always = BoolMatcher(true)

// Never is a Matcher which never matches
var Never = BoolMatcher(false)
