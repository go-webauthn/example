package middleware

import (
	"github.com/go-webauthn/example/internal/model"
	"go.uber.org/zap"
)

func Authenticated(next RequestHandler) (handler RequestHandler) {
	return func(ctx *RequestCtx) {
		session, err := ctx.GetUserSession()
		if err != nil || session.Username == "" {
			ctx.Log.Error("resource denied to anonymous user", zap.Error(err))

			ctx.ForbiddenJSON(model.NewErrorJSON().WithErrorStr("this is a private error").WithInfo("Forbidden."))

			return
		}

		next(ctx)
	}
}
