package middleware

import (
	"net/http"
)

// Middleware makes a type for all middleware functions so they can be 
// chained.
type Middleware func(http.Handler) http.handler

// Chain is meant to store a number of middleware that should be chained one
// after the other.
type Chain []Middleware

// Run chains together one or more Middleware, with the passed in
// http.HandlerFunc being run at the end.
func Run(h http.HandlerFunc, mws ...Middleware) http.Handler {
	// Reverse the middleware so they are run in the order they're passed in
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
