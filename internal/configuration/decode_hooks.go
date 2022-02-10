package configuration

import (
	"net/url"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap/zapcore"
)

func StringToZapCoreLevelHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		levelType := zapcore.InfoLevel

		switch t {
		case reflect.TypeOf(levelType):
			level, err := zapcore.ParseLevel(data.(string))
			if err != nil {
				return data, err
			}

			return level, nil
		case reflect.TypeOf(&levelType):
			level, err := zapcore.ParseLevel(data.(string))
			if err != nil {
				return data, err
			}

			return &level, nil
		default:
			return data, nil
		}
	}
}

func StringToURLHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		switch t {
		case reflect.TypeOf(url.URL{}):
			u, err := url.Parse(data.(string))
			if err != nil {
				return data, err
			}

			return *u, nil
		case reflect.TypeOf(&url.URL{}):
			u, err := url.Parse(data.(string))
			if err != nil {
				return data, err
			}

			return u, nil
		default:
			return data, nil
		}
	}
}
