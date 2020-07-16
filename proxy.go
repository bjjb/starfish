package starfish

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Proxy creates a new proxy (using the scheme) to the given host.
func Proxy(url *url.URL) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Host = url.Host
			r.URL.Scheme = url.Scheme
		},
	}
}
