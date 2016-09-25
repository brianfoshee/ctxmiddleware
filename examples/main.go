package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mw "github.com/brianfoshee/ctxmiddleware"
)

func main() {
	// A logger is used as an example of how to pass around some sort of object
	// to middleware and handlers. It may very well be a struct that contains a
	// logger, metrics collector, etc.
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	// In this case only a single middleware is inserted, timer(), but the
	// idea is for a slice of middleware to be passed as variadic args.
	// These would form a chain which passes along a context (now contained in
	// the *http.Request) as of Go 1.7.
	http.Handle("/", mw.Run(rootHandler(l), timer(l)))

	http.ListenAndServe(":8080", nil)
}

type key int

var userIPKey = 0

func timer(l *log.Logger) mw.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			now := time.Now()
			ip := req.RemoteAddr
			ctx := context.WithValue(req.Context(), userIPKey, ip)
			h.ServeHTTP(rw, req.WithContext(ctx))
			total := time.Since(now)
			l.Printf("path=%s, total=%.2fms", req.URL.Path, float64(total.Nanoseconds())/1e6)
		})
	}
}

func rootHandler(l *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, ok := r.Context().Value(userIPKey).(string)
		if !ok {
			l.Print("did not find ip in context")
			return
		}
		l.Print("in the root handler.")
		fmt.Fprintf(w, "Hello, ip is %s", ip)
	})
}
