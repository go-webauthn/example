package handler

import (
	"bytes"

	"github.com/go-webauthn/example/internal/middleware"
	"github.com/go-webauthn/example/internal/model"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"go.uber.org/zap"
)

func AssertionGET(ctx *middleware.RequestCtx) {
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

	var opts = []webauthn.LoginOption{
		webauthn.WithUserVerification(protocol.VerificationPreferred),
	}

	discoverable := ctx.QueryArgs().GetBool(queryArgDiscoverable)

	if !discoverable {
		if user, err = ctx.Providers.User.Get(session.Username); err != nil {
			ctx.Log.Error("failed to retrieve user", zap.Error(err))

			ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Forbidden."))

			return
		}

		credentials := user.WebAuthnCredentialDescriptors()

		opts = append(opts, webauthn.WithAllowedCredentials(credentials), webauthn.WithAppIdExtension(ctx.Config.ExternalURL.String()))
	}

	var (
		assertion *protocol.CredentialAssertion
		data      *webauthn.SessionData
	)

	if discoverable {
		ctx.Log.Debug("begin assertion", zap.Bool("discoverable", true))

		if assertion, data, err = ctx.Providers.Webauthn.BeginDiscoverableLogin(opts...); err != nil {
			ctx.Log.Error("error in begin discoverable assertion", append([]zap.Field{zap.Bool("discoverable", true)}, model.ProtoErrToFields(err)...)...)

			ctx.UnauthorizedJSON(model.NewErrorJSON().WithError(err).WithInfo("Unauthorized."))

			return
		}
	} else {
		ctx.Log.Debug("begin assertion", zap.Bool("discoverable", false), zap.String("user", user.ID))

		if assertion, data, err = ctx.Providers.Webauthn.BeginLogin(user, opts...); err != nil {
			ctx.Log.Error("error in begin assertion", append([]zap.Field{zap.Bool("discoverable", false), zap.String("user", user.ID)}, model.ProtoErrToFields(err)...)...)

			ctx.UnauthorizedJSON(model.NewErrorJSON().WithError(err).WithInfo("Unauthorized."))

			return
		}
	}

	session.Webauthn = data
	if err = ctx.SaveUserSession(session); err != nil {
		ctx.Log.Error("failed to save user session", zap.Error(err))

		ctx.UnauthorizedJSON(model.NewErrorJSON().WithError(err).WithInfo("Unauthorized."))

		return
	}

	ctx.OKJSON(assertion)
}

func AssertionPOST(ctx *middleware.RequestCtx) {
	var (
		user           *model.User
		credential     *webauthn.Credential
		parsedResponse *protocol.ParsedCredentialAssertionData
		session        *model.UserSession
		err            error
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

	discoverable := ctx.QueryArgs().GetBool(queryArgDiscoverable)

	if session == nil {

	}

	if parsedResponse, err = protocol.ParseCredentialRequestResponseBody(bytes.NewReader(ctx.PostBody())); err != nil {
		ctx.Log.Error("failed to parse credential request response body", model.ProtoErrToFields(err)...)

		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Forbidden."))

		return
	}

	if discoverable {
		if credential, err = ctx.Providers.Webauthn.ValidateDiscoverableLogin(func(_, userHandle []byte) (_ webauthn.User, err error) {
			if user, err = ctx.Providers.User.Get(string(userHandle)); err != nil {
				return nil, err
			}

			return user, nil
		}, *session.Webauthn, parsedResponse); err != nil {
			ctx.Log.Error("failed to validate discoverable assertion", model.ProtoErrToFields(err)...)

			ctx.UnauthorizedJSON(model.NewErrorJSON().WithError(err).WithInfo("Unauthorized."))

			return
		}

		session.Username = user.Name
	} else {
		if user, err = ctx.Providers.User.Get(session.Username); err != nil {
			ctx.Log.Error("failed to lookup user from session username", zap.String("username", session.Username), zap.Error(err))

			ctx.BadRequestJSON(model.NewErrorJSON().WithErrorStr("failed to lookup user").WithInfo("Bad Request."))

			return
		}

		if credential, err = ctx.Providers.Webauthn.ValidateLogin(user, *session.Webauthn, parsedResponse); err != nil {
			ctx.Log.Error("failed to validate assertion", model.ProtoErrToFields(err)...)

			ctx.UnauthorizedJSON(model.NewErrorJSON().WithError(err).WithInfo("Unauthorized."))

			return
		}
	}

	user.CredentialsSignIn = append(user.CredentialsSignIn, *credential)

	if err = ctx.Providers.User.Set(user); err != nil {
		ctx.Log.Error("failed to save user", zap.String("user_id", user.ID), zap.Error(err))

		ctx.ForbiddenJSON(model.NewErrorJSON().WithError(err).WithInfo("Forbidden."))

		return
	}

	ctx.OKJSON(user)
}
