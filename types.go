package webauthn

import (
	"github.com/fasthttp/session/v2"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/valyala/fasthttp"
)

type UserSession struct {
	Username string
	Webauthn *webauthn.SessionData
}

type User struct {
	ID                string
	Name              string                `json:"name"`
	DisplayName       string                `json:"display_name"`
	Icon              string                `json:"icon,omitempty"`
	Credentials       []webauthn.Credential `json:"credentials,omitempty"`
	CredentialsSignIn []webauthn.Credential `json:"credentials_sign_in,omitempty"`
}

func (u User) WebAuthnID() []byte {
	return []byte(u.ID)
}

func (u User) WebAuthnName() string {
	return u.Name
}

func (u User) WebAuthnDisplayName() string {
	return u.DisplayName
}

func (u User) WebAuthnIcon() string {
	return u.Icon
}

func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

type WebauthnCtx struct {
	fasthttp.RequestCtx

	Providers Providers
}

type Providers struct {
	User     UserProvider
	Session  *session.Session
	Webauthn *webauthn.WebAuthn
}

type UserProvider interface {
	GetUser(username string) (user *User, err error)
}
