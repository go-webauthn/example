package middleware

import (
	"github.com/valyala/fasthttp"
)

func CORS(next RequestHandler) RequestHandler {
	return func(ctx *RequestCtx) {
		accessControlRequestHeaders := ctx.Request.Header.Peek(fasthttp.HeaderAccessControlRequestHeaders)
		accessControlRequestMethod := ctx.Request.Header.Peek(fasthttp.HeaderAccessControlRequestMethod)

		if accessControlRequestHeaders != nil {
			ctx.Response.Header.SetBytesV(fasthttp.HeaderAccessControlAllowHeaders, accessControlRequestHeaders)
		}

		if accessControlRequestMethod != nil {
			ctx.Response.Header.SetBytesV(fasthttp.HeaderAccessControlAllowMethods, accessControlRequestMethod)
		}

		if origin := ctx.Request.Header.Peek(fasthttp.HeaderOrigin); origin != nil && (accessControlRequestMethod != nil || accessControlRequestHeaders != nil) {
			ctx.Response.Header.SetBytesV(fasthttp.HeaderAccessControlAllowOrigin, origin)
		}

		next(ctx)
	}
}
