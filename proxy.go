package starfish

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ProxyURL creates a new proxy (using the scheme) to the given host.
func ProxyURL(origin *url.URL) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.Header.Add("X-Forwarded-Host", r.Host)
			r.Header.Add("X-Origin-Host", origin.Host)
			r.URL.Host = origin.Host
			r.URL.Scheme = origin.Scheme
		},
	}
}

// Proxy creates a new proxy to the given host. If the URL is invalid, it'll
// serve a BadGateway instead.
func Proxy(upstream string) http.Handler {
	if url, err := url.Parse(upstream); err == nil {
		return ProxyURL(url)
	}
	return BadGateway
}
