// Vhost creates a lightweight router for routing http applications
// via. configured hostnames
package vhost

import (
	"net/http"
	"strings"
)

// Router
// allows registration of routes based on specific hostnames
//
// This satisfies the http.Handler interface so that it can be used to server requests
// Usage:
// 	net/http
//  	func main() {
//          router := vhost.New()
//			router.Handler(myHandler, "domain.com")
//      	http.Handle("/", router)
//      }
type Router struct {
	// we can assume that we will not need to concurrently modify routes
	routes map[string]*Route

	// Whether or not we will strictly match `www.` or not
	Strict bool

	// Configurable handler as a catch all
	NotFound http.Handler
}

// New
// creates a new hostname router
func New() *Router {
	return &Router{
		routes: make(map[string]*Route),
		Strict: false,
	}
}

// ServeHTTP
// satisfies the http.Handler requirements.
// when there are no routes to match we will fallback to the NotFound handler
func (r *Router) ServeHTTP(w http.ResponseWriter, rq *http.Request) {

	host := rq.URL.Host
	// due to mocking hostnames this block was added to ensure we can test properly
	if len(host) == 0 && len(rq.Host) != 0 {
		pos := strings.Split(rq.Host, ":")
		if len(pos) > 0 {
			host = pos[0]
		}
	}

	// No valid house
	if len(host) == 0 {
		r.fallback(w, rq)
		return
	}

	route := r.GetRoute(r.stripStrict(host))
	if route == nil {
		r.fallback(w, rq)
		return
	}

	route.Get().ServeHTTP(w, rq)
}

// GetRoute
// obtain an already registered route (if it exists)
func (r *Router) GetRoute(name string) *Route {
	return r.routes[name]
}

// fallback
// When no valid routes have been found we will fallback to a default `NotFound` response
func (r *Router) fallback(w http.ResponseWriter, rq *http.Request) {
	if r.NotFound == nil {
		http.NotFound(w, rq)
		return
	}

	r.NotFound.ServeHTTP(w, rq)
}

// Handler
// binds a http.Handler to the specified hostnames
func (r *Router) Handler(handler http.Handler, hosts ...string) *Router {

	for _, host := range hosts {

		r.routes[host] = &Route{
			name:    r.stripStrict(host),
			handler: handler,
		}
	}

	return r
}

// HandlerFunc
// binds http.HandlerFunc to the specific hostnames
func (r *Router) HandlerFunc(f func(w http.ResponseWriter, r *http.Request), hosts ...string) *Router {
	return r.Handler(http.HandlerFunc(f), hosts...)
}

// stripStrict
// performs actions to strip a host based on the current `strict` ruleset
func (r *Router) stripStrict(host string) string {
	// remove any possible www. when not in strict mode
	if !r.Strict && host[0:4] == "www." {
		host = host[4:]
	}

	return host
}
