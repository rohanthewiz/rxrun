package main

import (
	"fmt"
	"github.com/rohanthewiz/rxrouter"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	rx := rxrouter.New(
		rxrouter.Options{Verbose: true, Port: "3026"}, // the Options argument here is optional
	)

	// Logging middleware
	rx.Use(
		func(ctx *fasthttp.RequestCtx) (ok bool) {
			log.Printf("Requested path: %s", ctx.Path())
			return true
		},
		fasthttp.StatusServiceUnavailable, // 503
	)

	// Auth middleware
	rx.Use(
		func(ctx *fasthttp.RequestCtx) (ok bool) {
			authed := true // pretend we got a good response from our auth check
			if !authed {
				return false
			}
			return true
		},
		fasthttp.StatusUnauthorized,
	)

	// Add some routes
	rx.AddRoute("/", func(ctx *fasthttp.RequestCtx, params map[string]string) {
		fmt.Fprintf(ctx, "Hello, world! Requested path is %q", ctx.Path())
	})
	rx.AddRoute("/hello/:name/:age", handleHello)
	rx.AddRoute("/store/:number/:location", handleStore)
	// Routes for static files
	rx.AddStaticFilesRoute("/images/", "./assets/images", 1)
	rx.AddStaticFilesRoute("/css/", "./assets/css", 1)

	// Let it rip!
	rx.Start()
}

func handleHello(ctx *fasthttp.RequestCtx, params map[string]string) {
	ctx.WriteString(fmt.Sprintf("Hello %s", params["name"]))
}

func handleStore(ctx *fasthttp.RequestCtx, params map[string]string) {
	fmt.Fprintf(ctx, "Store: %s  location: %s", params["number"], params["location"])
}
