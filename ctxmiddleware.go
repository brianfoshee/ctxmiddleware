// Package ctxmiddleware should not be used anymore now that x/net/context has
// been merged into the stdlib as of go 1.7. This is a basic middleware
// chaining implementation.
// See old versions of this file for go versions previous to 1.7.
package ctxmiddleware

import "net/http"

// Middleware makes a type for all middleware functions so they can be
// chained.
type Middleware func(http.Handler) http.Handler

// Chain is meant to store a number of middleware that should be chained one
// after the other.
type Chain []Middleware

// Run chains together one or more Middleware, with the passed in
// http.Handler being run at the end.
func Run(h http.Handler, mws ...Middleware) http.Handler {
	// Reverse the middleware so they are run in the order they're passed in
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
