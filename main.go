package main

import (
	"fmt"
	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/rxrouter"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"strings"
)

func main() {
	rx := rxrouter.New(rxrouter.Options{
		Verbose: true,
		TLS: rxrouter.RxTLS{
			UseTLS:   false,
			CertFile: "/etc/letsencrypt/live/mysite.com/cert.pem",
			KeyFile:  "/etc/letsencrypt/live/mysite.com/privkey.pem",
		},
	})
	// Env port override
	prt := os.Getenv("RX_PORT")
	if prt != "" {
		rx.Options.Port = prt
	}
	tls := strings.ToLower(os.Getenv("RX_USE_TLS"))
	if tls != "" && tls != "false" && tls != "no" {
		rx.Options.TLS.UseTLS = true
		go startHttpToHTTPSRedirector()
	}

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
	rx.AddRoute("/", handleRoot)
	rx.AddRoute("/hello/:name/:age", func(ctx *fasthttp.RequestCtx, params map[string]string) {
		ctx.WriteString(fmt.Sprintf("Hello %s. You are %s!", params["name"], params["age"]))
	})
	// Routes for static files
	rx.AddStaticFilesRoute("/images/", "./assets/images", 1)
	rx.AddStaticFilesRoute("/css/", "./assets/css", 1)
	rx.AddStaticFilesRoute("/.well-known/acme-challenge/", "./certs", 0)

	// Let it rip!
	rx.Start()
}

func handleRoot(ctx *fasthttp.RequestCtx, params map[string]string) {
	e := element.New
	htm := "<!DOCTYPE html>" +
		e("html").R(
			e("head").R(
				e("title").R("Under Maintenance"),
			),
			e("body").R(
				e("div", "class", "main-content").R(
					e("h2").R("This site is currently under maintenance. Please check back later, and thanks for your patience"),
				),
			),
		)
	ctx.SetContentType("text/html; charset=utf-8")
	ctx.WriteString(htm)
}

// Use a simple router instance to redirect traffic from standard to TLS port
// Typically this is 80 -> 443
func startHttpToHTTPSRedirector() {
	fx := rxrouter.New(rxrouter.Options{
		Verbose: true,
		Port:    "80",
		TLS: rxrouter.RxTLS{
			UseTLS: false, // just being explicit here
		},
	})

	fx.AddCustomMasterHandler(func(ctx *fasthttp.RequestCtx) {
		ctx.Redirect("https://mysite.com", fasthttp.StatusTemporaryRedirect)
	})

	fx.Start()
}
