// Package ctxmiddleware is an example of using custom context handlers and
// middleware. It is largely based on the final option in this article:
// https://joeshaw.org/net-context-and-http-handler/.
// I have added a middleware chaining method and an associated ContextMW type.
//
// The idea is that in the case of this proposal being accepted,
// https://github.com/golang/go/issues/14660#issuecomment-193914014,
// where the context package would be in the standard library and implemented
// on http.Request, that these same handlers and middleware could be used after
// a brief refactoring to remove the custom types.
package ctxmiddleware

import (
	"net/http"

	"golang.org/x/net/context"
)

// ContextHandler models itself off http.Handler with the addition of a
// Context object.
type ContextHandler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request)
}

// ContextHandlerFunc models itself off http.HandlerFunc with the addition of a
// Context object.
type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

// ServeHTTPContext makes ContextHandlerFunc adhere to ContextHandler
func (h ContextHandlerFunc) ServeHTTPContext(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	h(ctx, rw, req)
}

// ContextAdapter pairs a context object with an http handler. This can be
// passed into http.Handle as it implements ServeHTTP.
type ContextAdapter struct {
	Context context.Context
	Handler ContextHandler
}

// ServeHTTP causes ContextAdapter to adhere to the http.Handler interface
func (ca *ContextAdapter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ca.Handler.ServeHTTPContext(ca.Context, rw, req)
}

// ContextMW makes a type for all middleware functions so they can be chained.
type ContextMW func(ContextHandler) ContextHandler

// Run chains together one or more ContextMW, with the passed in Contexthandler
// being run at the end.
func Run(h ContextHandler, mws ...ContextMW) ContextHandler {
	// Reverse the middleware so they are run in the order they're passed in
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
