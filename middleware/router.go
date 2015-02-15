package middleware

import (
	"net/http"
	"strings"

	"github.com/casualjim/go-swagger/errors"
	"github.com/gorilla/context"
)

func newRouter(ctx *Context, next http.Handler) http.Handler {
	isRoot := ctx.spec.BasePath() == "" || ctx.spec.BasePath() == "/"

	handleRequest := func(rw http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(rw, r)
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer context.Clear(r)
		// use context to lookup routes
		if isRoot {
			if _, ok := ctx.RouteInfo(r); ok {
				handleRequest(rw, r)
				return
			}
		} else {
			if p := strings.TrimPrefix(r.URL.Path, ctx.spec.BasePath()); len(p) < len(r.URL.Path) {
				r.URL.Path = p
				if _, ok := ctx.RouteInfo(r); ok {
					handleRequest(rw, r)
					return
				}
			}
		}
		// Not found, check if it exists in the other methods first
		if others := ctx.AllowedMethods(r); len(others) > 0 {
			ctx.Respond(rw, r, ctx.spec.RequiredProduces(), nil, errors.MethodNotAllowed(r.Method, others))
			return
		}

		ctx.Respond(rw, r, ctx.spec.RequiredProduces(), nil, errors.NotFound("path %s was not found", r.URL.Path))
	})
}