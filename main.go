// This package is an example of using custom context handlers and middleware.
// It is largely based on the final option in this article:
// https://joeshaw.org/net-context-and-http-handler/
// I have added a middleware chaining method and an associated ContextMW type.
//
// The idea is that in the case of this proposal being accepted,
// https://github.com/golang/go/issues/14660#issuecomment-193914014,
// where the context package would be in the standard library and implemented
// on http.Request, that these same handlers and middleware could be used after
// a brief refactoring to remove the custom types.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
)

func main() {
	// A logger is used as an example of how to pass around some sort of object
	// to middleware and handlers. It may very well be a struct that contains a
	// logger, metrics collector, etc.
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	h := &ContextAdapter{
		ctx: context.Background(),

		// In this case only a single middleware is inserted, timer(), but the
		// idea is for a slice of middleware to be passed as variadic args.
		// These would form a chain which passes along a context.
		handler: Run(rootHandler(l), timer(l)),
	}

	http.Handle("/", h)
	http.ListenAndServe(":8080", nil)
}

func timer(l *log.Logger) ContextMW {
	return func(h ContextHandler) ContextHandler {
		return ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			now := time.Now()
			h.ServeHTTPContext(ctx, rw, req)
			end := time.Now()
			total := end.Sub(now)
			l.Printf("path=%s, total=%.2fms", req.URL.Path, float64(total.Nanoseconds())/1e6)
		})
	}
}

func rootHandler(l *log.Logger) ContextHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		l.Print("in the root handler")
		fmt.Fprintf(w, "Hello")
	}
}

type ContextHandler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request)
}

type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (h ContextHandlerFunc) ServeHTTPContext(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	h(ctx, rw, req)
}

// ContextHandler adheres to http.Handler interface
type ContextAdapter struct {
	ctx     context.Context
	handler ContextHandler
}

// ServeHTTP causes ContextAdapter to adhere to http.Handler interface
func (ca *ContextAdapter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ca.handler.ServeHTTPContext(ca.ctx, rw, req)
}

type ContextMW func(ContextHandler) ContextHandler

func Run(h ContextHandler, mws ...ContextMW) ContextHandler {
	// these will run last passed-in first
	for _, mw := range mws {
		fmt.Println("adding mw")
		h = mw(h)
	}
	return h
}
