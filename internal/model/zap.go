package model

import (
	"github.com/go-webauthn/webauthn/protocol"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ProtoErrToFields(err error) (fields []zap.Field) {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *protocol.Error:
		return []zap.Field{
			{Key: "err", Type: zapcore.ErrorType, Interface: e},
			{Key: "details", Type: zapcore.StringType, String: e.Details},
			{Key: "info", Type: zapcore.StringType, String: e.DevInfo},
			{Key: "type", Type: zapcore.StringType, String: e.Type},
		}
	default:
		return nil
	}
}
