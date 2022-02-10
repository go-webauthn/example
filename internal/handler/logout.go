package handler

import (
	"github.com/go-webauthn/example/internal/middleware"
	"github.com/go-webauthn/example/internal/model"
)

func LogoutGET(ctx *middleware.RequestCtx) {
	if err := ctx.DestroyUserSession(); err != nil {
		ctx.ForbiddenJSON(model.NewErrorJSON().WithErrorStr("could not regenerate session").WithInfo("Invalid Credentials."))

		return
	}

	ctx.OKJSON(nil)
}
