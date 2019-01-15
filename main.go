package main

import (
	"fmt"
	"github.com/rohanthewiz/rxrouter"
	"github.com/rohanthewiz/rxrouter/mux"
	"github.com/valyala/fasthttp"
	"log"
)

const appEnv = "[DEV]"
const authed = true

func main() {
	rx := rxrouter.New()

	// Rudimentary request logging middleware
	rx.Use(func(ctx *fasthttp.RequestCtx) (retCtx *fasthttp.RequestCtx, ok bool) {
		log.Printf("Requested path: %s", ctx.Path())
		return ctx, true
	}, fasthttp.StatusServiceUnavailable) // 503

	// Auth middleware
	rx.Use(func(ctx *fasthttp.RequestCtx) (retCtx *fasthttp.RequestCtx, ok bool) {
		if !authed { return ctx, false }
		return ctx, true
	}, fasthttp.StatusUnauthorized)

	// Prepend to output middleware
	rx.Use(func(ctx *fasthttp.RequestCtx) (retCtx *fasthttp.RequestCtx, ok bool) {
		_, _ = fmt.Fprintf(ctx, "%s ", appEnv)
		return ctx, true
	}, fasthttp.StatusNotImplemented)


	// Add some routes
	rx.Mux.Add("/", func (ctx *fasthttp.RequestCtx, mx *mux.Mux) {
		_, _ = fmt.Fprintf(ctx, "Hello, world! Requested path is %q", ctx.Path())
	})
	rx.Mux.Add("/abc", handleABC)


	// Let it rip!
	rx.Start("3020")
}
func handleABC(ctx *fasthttp.RequestCtx, mx *mux.Mux) {
	_, _ = fmt.Fprintf(ctx, "Hello ABC! Requested path is %q", ctx.Path())
}
