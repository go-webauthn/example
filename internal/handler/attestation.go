package handler

import (
	"bytes"

	"github.com/go-webauthn/example/internal/middleware"
	"github.com/go-webauthn/example/internal/model"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"go.uber.org/zap"
)

func AttestationGET(ctx *middleware.RequestCtx) {
	var (
		user    *model.User
		session *model.UserSession
		err     error
	)

	if session, err = ctx.GetUserSession(); err != nil {
		ctx.Log.Error("failed to retrieve user session", zap.Error(err))

		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Forbidden."))

		return
	}

	if user, err = ctx.Providers.User.Get(session.Username); err != nil {
		ctx.Log.Error("failed to retrieve user", zap.Error(err))

		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Forbidden."))

		return
	}

	discoverable := ctx.QueryArgs().GetBool(queryArgDiscoverable)

	var selection protocol.AuthenticatorSelection

	if discoverable {
		selection = ctx.Config.AuthenticatorSelection(protocol.ResidentKeyRequirementRequired)
	} else {
		selection = ctx.Config.AuthenticatorSelection(protocol.ResidentKeyRequirementDiscouraged)
	}

	opts, data, err := ctx.Providers.Webauthn.BeginRegistration(user,
		webauthn.WithAuthenticatorSelection(selection),
		webauthn.WithConveyancePreference(ctx.Config.ConveyancePreference),
		webauthn.WithExclusions(user.WebAuthnCredentialDescriptors()),
		webauthn.WithAppIdExcludeExtension(ctx.Config.ExternalURL.String()),
	)

	if err != nil {
		ctx.Log.Error("failed to generate attestation options", model.ProtoErrToFields(err)...)

		ctx.UnauthorizedJSON(model.NewErrorJSON().WithError(err).WithInfo("Unauthorized."))

		return
	}

	session.Webauthn = data

	if err = ctx.SaveUserSession(session); err != nil {
		ctx.Log.Error("failed to save user session", zap.String("username", session.Username), zap.Error(err))

		ctx.UnauthorizedJSON(model.NewErrorJSON().WithError(err).WithInfo("Unauthorized."))

		return
	}

	ctx.OKJSON(opts)
}

func AttestationPOST(ctx *middleware.RequestCtx) {
	var (
		user    *model.User
		session *model.UserSession
		err     error
	)

	if session, err = ctx.GetUserSession(); err != nil {
		ctx.Log.Error("failed to retrieve user session", zap.Error(err))

		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Forbidden."))

		return
	}

	defer func() {
		session.Webauthn = nil

		if err := ctx.SaveUserSession(session); err != nil {
			ctx.Log.Error("failed to save user session", zap.Error(err))
		}
	}()

	if user, err = ctx.Providers.User.Get(session.Username); err != nil {
		ctx.Log.Error("failed to retrieve user", zap.Error(err))

		ctx.UnauthorizedJSON(model.NewErrorJSON().WithError(err).WithInfo("Unauthorized."))

		return
	}

	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(ctx.PostBody()))
	if err != nil {
		ctx.Log.Error("failed to parse credential creation response body", model.ProtoErrToFields(err)...)

		ctx.BadRequestJSON(model.NewErrorJSON().WithError(err).WithInfo("Bad Request."))

		return
	}

	cred, err := ctx.Providers.Webauthn.CreateCredential(user, *session.Webauthn, parsedResponse)
	if err != nil {
		ctx.Log.Error("failed to create credential", model.ProtoErrToFields(err)...)

		ctx.UnauthorizedJSON(model.NewErrorJSON().WithError(err).WithInfo("Unauthorized."))

		return
	}

	user.Credentials = append(user.Credentials, *cred)

	if err = ctx.Providers.User.Set(user); err != nil {
		ctx.Log.Error("failed to save user", zap.String("user_id", user.ID), zap.Error(err))

		ctx.UnauthorizedJSON(model.NewErrorJSON().WithError(err).WithInfo("Unauthorized."))

		return
	}

	ctx.CreatedJSON("Done.")
}
