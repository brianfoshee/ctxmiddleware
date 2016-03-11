package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mw "github.com/brianfoshee/ctxmiddleware"

	"golang.org/x/net/context"
)

func main() {
	// A logger is used as an example of how to pass around some sort of object
	// to middleware and handlers. It may very well be a struct that contains a
	// logger, metrics collector, etc.
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	h := &mw.ContextAdapter{
		Context: context.Background(),

		// In this case only a single middleware is inserted, timer(), but the
		// idea is for a slice of middleware to be passed as variadic args.
		// These would form a chain which passes along a context.
		Handler: mw.Run(rootHandler(l), timer(l)),
	}

	http.Handle("/", h)
	http.ListenAndServe(":8080", nil)
}

func timer(l *log.Logger) mw.ContextMiddleware {
	return func(h mw.ContextHandler) mw.ContextHandler {
		return mw.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			now := time.Now()
			h.ServeHTTPContext(ctx, rw, req)
			end := time.Now()
			total := end.Sub(now)
			l.Printf("path=%s, total=%.2fms", req.URL.Path, float64(total.Nanoseconds())/1e6)
		})
	}
}

func rootHandler(l *log.Logger) mw.ContextHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		l.Print("in the root handler")
		fmt.Fprintf(w, "Hello")
	}
}
