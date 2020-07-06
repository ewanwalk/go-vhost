package vhost

import "net/http"

type Route struct {
	name    string
	handler http.Handler
}

// Handler
// sets a handler
func (r *Route) Handler(handler http.Handler) *Route {
	r.handler = handler
	return r
}

// HandlerFunc
// sets a handler function
func (r *Route) HandlerFunc(f func(response http.ResponseWriter, r *http.Request)) *Route {
	return r.Handler(http.HandlerFunc(f))
}

// Get
// Returns the http handler
func (r *Route) Get() http.Handler {
	return r.handler
}
