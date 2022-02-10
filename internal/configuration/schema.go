package configuration

import (
	"net/url"

	"github.com/go-webauthn/webauthn/protocol"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Log Log `koanf:"log"`

	ListenAddress string  `koanf:"listen_address"`
	ExternalURL   url.URL `koanf:"external_url"`
	Session       Session `koanf:"session"`
	DisplayName   string  `koanf:"display_name"`

	UserVerification        protocol.UserVerificationRequirement `koanf:"user_verification_requirement"`
	AuthenticatorAttachment protocol.AuthenticatorAttachment     `koanf:"authenticator_attachment"`
	ConveyancePreference    protocol.ConveyancePreference        `koanf:"conveyance_preference"`
}

func (c Config) AuthenticatorSelection(requirement protocol.ResidentKeyRequirement) (selection protocol.AuthenticatorSelection) {
	selection = protocol.AuthenticatorSelection{
		AuthenticatorAttachment: c.AuthenticatorAttachment,
		UserVerification:        c.UserVerification,
		ResidentKey:             requirement,
	}

	if selection.ResidentKey == "" {
		selection.ResidentKey = protocol.ResidentKeyRequirementDiscouraged
	}

	switch selection.ResidentKey {
	case protocol.ResidentKeyRequirementRequired:
		selection.RequireResidentKey = protocol.ResidentKeyRequired()
	case protocol.ResidentKeyRequirementDiscouraged:
		selection.RequireResidentKey = protocol.ResidentKeyNotRequired()
	}

	if selection.AuthenticatorAttachment == "" {
		selection.AuthenticatorAttachment = protocol.CrossPlatform
	}

	if selection.UserVerification == "" {
		selection.UserVerification = protocol.VerificationPreferred
	}

	return selection
}

type Log struct {
	Level zapcore.Level `koanf:"level"`

	File    FileLog    `koanf:"file"`
	Console ConsoleLog `koanf:"console"`
}

type FileLog struct {
	Level    *zapcore.Level `koanf:"level"`
	Encoding string         `koanf:"encoding"`
	Path     string         `koanf:"path"`
}

type ConsoleLog struct {
	Level    *zapcore.Level `koanf:"level"`
	Encoding string         `koanf:"encoding"`
	Disable  bool           `koanf:"disable"`
}

type Session struct {
	CookieName string `koanf:"cookie_name"`
	Secure     bool   `koanf:"secure"`
	Domain     string `koanf:"domain"`
}
