package starfish

import "net/http"

func New() http.Handler {
	return &Router{}
}
