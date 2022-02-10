package handler

import (
	"crypto/elliptic"
	"encoding/json"

	"github.com/go-webauthn/example/internal/middleware"
	"github.com/go-webauthn/example/internal/model"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/tstranex/u2f"
	"go.uber.org/zap"
)

func U2FRegisterGET(ctx *middleware.RequestCtx) {
	appID := ctx.Config.ExternalURL.String()

	challenge, err := u2f.NewChallenge(appID, []string{appID})
	if err != nil {
		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Operation Failed."))

		return
	}

	session, err := ctx.GetUserSession()
	if err != nil {
		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Operation Failed."))

		return
	}

	user, err := ctx.Providers.User.Get(session.Username)
	if err != nil {
		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Operation Failed."))

		return
	}

	regs := user.U2FRegistrations()

	session.U2F = challenge

	if err = ctx.SaveUserSession(session); err != nil {
		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Operation Failed."))

		return
	}

	ctx.OKJSON(u2f.NewWebRegisterRequest(challenge, regs))
}

func U2FRegisterPOST(ctx *middleware.RequestCtx) {
	registerResponse := u2f.RegisterResponse{}

	if err := json.Unmarshal(ctx.PostBody(), &registerResponse); err != nil {
		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Operation Failed."))

		return
	}

	session, err := ctx.GetUserSession()
	if err != nil {
		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Operation Failed."))

		return
	}

	if session.U2F == nil {
		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Operation Failed."))

		return
	}

	defer func() {
		session.U2F = nil

		if err := ctx.SaveUserSession(session); err != nil {
			ctx.Log.Error("failed to save user session", zap.Error(err))
		}
	}()

	registration, err := u2f.Register(registerResponse, *session.U2F, &u2f.Config{SkipAttestationVerify: true})
	if err != nil {
		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Operation Failed."))

		return
	}

	user, err := ctx.Providers.User.Get(session.Username)
	if err != nil {
		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Operation Failed."))

		return
	}

	user.Credentials = append(user.Credentials, webauthn.Credential{
		AttestationType: model.WebauthnAttestationTypeFIDOU2F,
		ID:              registration.KeyHandle,
		PublicKey:       elliptic.Marshal(elliptic.P256(), registration.PubKey.X, registration.PubKey.Y),
	})

	err = ctx.Providers.User.Set(user)
	if err != nil {
		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Operation Failed."))

		return
	}

	ctx.CreatedJSON("Created.")
}
