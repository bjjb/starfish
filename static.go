package starfish

import (
	"net/http"
)

// Static creates a new http.FileServer for the given directory.
func Static(path string) http.Handler {
	return http.FileServer(http.Dir(path))
}
