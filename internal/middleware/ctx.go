package middleware

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fasthttp/session/v2"
	"github.com/fasthttp/session/v2/providers/memory"
	"github.com/go-webauthn/example/internal/model"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/go-webauthn/example/internal/configuration"
)

func NewRequestHandlerCtxMiddleware(config *configuration.Config, providers *Providers) (bridge RequestHandlerMiddleware) {
	return func(next RequestHandler) (handler fasthttp.RequestHandler) {
		return func(requestCtx *fasthttp.RequestCtx) {
			ctx := NewRequestCtx(requestCtx, config, providers)

			next(ctx)
		}
	}
}

func NewRequestCtx(requestCtx *fasthttp.RequestCtx, config *configuration.Config, providers *Providers) (ctx *RequestCtx) {
	ctx = new(RequestCtx)

	ctx.RequestCtx = requestCtx
	ctx.Config = config
	ctx.Providers = providers

	requestUUID, err := uuid.NewUUID()
	if err == nil {
		//ctx.Log = zap.L().With(zap.String("request_uuid", requestUUID.String()), zap.String("remote_ip", requestCtx.RemoteIP().String()))
		ctx.Log = zap.L().WithOptions(zap.Fields(zap.String("request_uuid", requestUUID.String()), zap.String("remote_ip", ctx.RemoteIP().String()), zap.ByteString("path", ctx.Path()), zap.ByteString("method", ctx.Method())))
	} else {
		ctx.Log = zap.L()
	}

	return ctx
}

func NewProviders(config *configuration.Config) (providers *Providers, err error) {
	providers = new(Providers)

	providers.User = NewMemoryUserProvider()

	sessionConfig := session.NewDefaultConfig()

	sessionConfig.CookieName = config.Session.CookieName
	if config.Session.Domain != "" && strings.HasSuffix(config.ExternalURL.Hostname(), config.Session.Domain) {
		sessionConfig.Domain = config.Session.Domain
	} else {
		sessionConfig.Domain = config.ExternalURL.Hostname()
	}

	sessionConfig.Secure = true
	sessionConfig.CookieSameSite = fasthttp.CookieSameSiteLaxMode

	providers.Session = session.New(sessionConfig)

	if sessionStoreProvider, err := memory.New(memory.Config{}); err != nil {
		return nil, err
	} else if err = providers.Session.SetProvider(sessionStoreProvider); err != nil {
		return nil, err
	}

	if providers.Webauthn, err = webauthn.New(&webauthn.Config{
		RPID:                  config.ExternalURL.Hostname(),
		RPDisplayName:         config.DisplayName,
		RPOrigin:              config.ExternalURL.String(),
		AttestationPreference: config.ConveyancePreference,
	}); err != nil {
		return nil, err
	}

	return providers, nil
}

func (ctx *RequestCtx) DestroyUserSession() (err error) {
	if err = ctx.RegenerateUserSession(); err != nil {
		return err
	}

	store := session.NewStore()

	return ctx.Providers.Session.Save(ctx.RequestCtx, store)
}

func (ctx *RequestCtx) RegenerateUserSession() (err error) {
	return ctx.Providers.Session.Regenerate(ctx.RequestCtx)
}

func (ctx *RequestCtx) GetUserSession() (session *model.UserSession, err error) {
	store, err := ctx.Providers.Session.Get(ctx.RequestCtx)
	if err != nil {
		return nil, err
	}

	session = &model.UserSession{}

	sessionBytes, ok := store.Get("user").([]byte)
	if !ok {
		return session, nil
	}

	if err = json.Unmarshal(sessionBytes, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (ctx *RequestCtx) SaveUserSession(session *model.UserSession) (err error) {
	store, err := ctx.Providers.Session.Get(ctx.RequestCtx)
	if err != nil {
		return err
	}

	sessionJSON, err := json.Marshal(*session)
	if err != nil {
		return err
	}

	store.Set("user", sessionJSON)

	return ctx.Providers.Session.Save(ctx.RequestCtx, store)
}

func (ctx *RequestCtx) CreateKO(message interface{}) (ko model.MessageResponse) {
	ko = model.MessageResponse{
		Status: "KO",
	}

	switch m := message.(type) {
	case *model.ErrorJSON:
		ko.Message = m.Info()
	case error:
		ko.Message = m.Error()
	case string:
		ko.Message = m
	}

	return ko
}

func (ctx *RequestCtx) CreatedJSON(message string) {
	response := model.MessageResponse{
		Status:  "OK",
		Message: message,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		ctx.Log.Error("failed to marshal JSON 201 Created response", zap.Error(err))

		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody(responseJSON)
}

func (ctx *RequestCtx) OKJSON(data interface{}) {
	response := model.DataResponse{
		Status: "OK",
		Data:   data,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		ctx.Log.Error("failed to marshal JSON 200 OK response", zap.Error(err))

		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(responseJSON)
}

func (ctx *RequestCtx) ErrorJSON(err error, status int) {
	response := ctx.CreateKO(err)

	responseJSON, err := json.Marshal(response)
	if err != nil {
		ctx.Log.Error(fmt.Sprintf("failed to marshal JSON %d %s response", status, fasthttp.StatusMessage(status)), zap.Error(err))

		return
	}

	ctx.SetStatusCode(status)
	ctx.SetBody(responseJSON)
}

func (ctx *RequestCtx) BadRequestJSON(err error) {
	ctx.ErrorJSON(err, fasthttp.StatusBadRequest)
}

func (ctx *RequestCtx) ForbiddenJSON(err error) {
	ctx.ErrorJSON(err, fasthttp.StatusForbidden)
}

func (ctx *RequestCtx) UnauthorizedJSON(err error) {
	ctx.ErrorJSON(err, fasthttp.StatusUnauthorized)
}
