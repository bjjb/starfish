package main

import "net/http"

type api string

func (api api) Match(r *http.Request) bool {
	return true
}

func (api api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusTeapot), http.StatusTeapot)
}
