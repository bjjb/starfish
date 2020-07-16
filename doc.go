// Package starfish provides a library for building a HTTP edge router. For
// every request, it will check against a list of matching functions, and
// serve the request with the appropriate handler. This can be used to build
// reverse proxies, file servers, edge gateways, etc - see the examples for
// inspiration.
package starfish
