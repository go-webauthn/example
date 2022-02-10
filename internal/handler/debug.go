package handler

import (
	"github.com/go-webauthn/example/internal/middleware"
)

func DebugGET(ctx *middleware.RequestCtx) {

	ctx.SetBodyString("Hello")
}
