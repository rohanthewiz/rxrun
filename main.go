package main

import (
	"fmt"
	"github.com/rohanthewiz/rxrouter"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	rx := rxrouter.New(rxrouter.Options{Verbose: true}) // the argument here is optional

	// Logging middleware
	rx.Use(func(ctx *fasthttp.RequestCtx) (ok bool) {
		log.Printf("Requested path: %s", ctx.Path())
		return true
	}, fasthttp.StatusServiceUnavailable) // 503

	// Auth middleware
	rx.Use(func(ctx *fasthttp.RequestCtx) (ok bool) {
		const authed = true // pretend we got a good response from our auth check
		if !authed {
			return false
		}
		return true
	}, fasthttp.StatusUnauthorized)

	// Add some routes
	rx.Add("/", func(ctx *fasthttp.RequestCtx, params map[string]string) {
		_, _ = fmt.Fprintf(ctx, "Hello, world! Requested path is %q", ctx.Path())
	})

	rx.Add("/hello/:name", handleHello)

	rx.Add("/images/:name/:height", handleImages)

	// Let it rip!
	rx.Start("3020")
}

func handleHello(ctx *fasthttp.RequestCtx, params map[string]string) {
	_, _ = ctx.WriteString(fmt.Sprintf("Hello %s", params["name"]))
}

func handleImages(ctx *fasthttp.RequestCtx, params map[string]string) {
	_, _ = fmt.Fprintf(ctx, "Image: %s  height: %s", params["name"], params["height"])
}
