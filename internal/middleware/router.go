package middleware

import (
	"github.com/fasthttp/router"
	"github.com/go-webauthn/example/internal/configuration"
	"github.com/valyala/fasthttp"
)

func NewRouter(config *configuration.Config, providers *Providers, middlewares ...Middleware) (r *Router) {
	r = &Router{
		middleware:  NewRequestHandlerCtxMiddleware(config, providers),
		middlewares: middlewares,
		router:      router.New(),
	}

	return r
}

type Router struct {
	middleware  RequestHandlerMiddleware
	middlewares []Middleware
	router      *router.Router
}

func (r *Router) wrap(handler RequestHandler) (middleware fasthttp.RequestHandler) {
	if len(r.middlewares) == 0 {
		return r.middleware(handler)
	}

	var h RequestHandler

	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}

	return r.middleware(h)
}

func (r *Router) Router() (router *router.Router) {
	return r.router
}

func (r *Router) Handle(method string, path string, handler RequestHandler) {
	r.router.Handle(method, path, r.wrap(handler))
}

func (r *Router) ANY(path string, handler RequestHandler) {
	r.router.ANY(path, r.wrap(handler))
}

func (r *Router) GET(path string, handler RequestHandler) {
	r.router.GET(path, r.wrap(handler))
}

func (r *Router) POST(path string, handler RequestHandler) {
	r.router.POST(path, r.wrap(handler))
}

func (r *Router) DELETE(path string, handler RequestHandler) {
	r.router.DELETE(path, r.wrap(handler))
}

func (r *Router) PUT(path string, handler RequestHandler) {
	r.router.PUT(path, r.wrap(handler))
}

func (r *Router) HEAD(path string, handler RequestHandler) {
	r.router.HEAD(path, r.wrap(handler))
}

func (r *Router) CONNECT(path string, handler RequestHandler) {
	r.router.CONNECT(path, r.wrap(handler))
}

func (r *Router) PATCH(path string, handler RequestHandler) {
	r.router.PATCH(path, r.wrap(handler))
}

func (r *Router) TRACE(path string, handler RequestHandler) {
	r.router.TRACE(path, r.wrap(handler))
}

func (r *Router) OPTIONS(path string, handler RequestHandler) {
	r.router.OPTIONS(path, r.wrap(handler))
}

func (r *Router) Handler(ctx *fasthttp.RequestCtx) {
	r.router.Handler(ctx)
}
