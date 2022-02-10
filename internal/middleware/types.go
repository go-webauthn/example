package middleware

import (
	"github.com/fasthttp/session/v2"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/go-webauthn/example/internal/configuration"
	"github.com/go-webauthn/example/internal/model"
)

type RequestCtx struct {
	*fasthttp.RequestCtx

	Log       *zap.Logger
	Providers *Providers
	Config    *configuration.Config
}

type UserProvider interface {
	Get(name string) (user *model.User, err error)
	Set(user *model.User) (err error)
}

type Providers struct {
	Webauthn *webauthn.WebAuthn
	Session  *session.Session

	User UserProvider
}

type Middleware func(next RequestHandler) (handler RequestHandler)

type RequestHandlerMiddleware = func(handler RequestHandler) fasthttp.RequestHandler

type RequestHandler = func(ctx *RequestCtx)
