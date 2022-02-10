package handler

import (
	"github.com/go-webauthn/example/internal/middleware"
	"github.com/go-webauthn/example/internal/model"
)

func InfoGET(ctx *middleware.RequestCtx) {
	session, err := ctx.GetUserSession()
	if err != nil {
		ctx.BadRequestJSON(model.NewErrorJSON().WithError(err).WithInfo("Could not get session."))
	}

	response := model.Info{
		Username: session.Username,
	}

	ctx.OKJSON(response)
}
