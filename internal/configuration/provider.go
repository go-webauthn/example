package configuration

import (
	"github.com/knadh/koanf"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
)

func Load(paths []string, environment bool, flags *pflag.FlagSet) (config *Config, err error) {
	config = &Config{}

	ko := koanf.NewWithConf(koanf.Conf{
		Delim:       StandardDelimiter,
		StrictMerge: false,
	})

	if len(paths) != 0 {
		if err = loadPaths(ko, paths, StandardDelimiter, false); err != nil {
			return nil, err
		}
	}

	if environment {
		if err = loadEnvironment(ko, EnvironmentDelimiter, false); err != nil {
			return nil, err
		}
	}

	if flags != nil {
		if err = loadFlags(ko, StandardDelimiter, false, flags); err != nil {
			return nil, err
		}
	}

	decodeHook := mapstructure.ComposeDecodeHookFunc(
		StringToURLHookFunc(),
		StringToZapCoreLevelHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
	)

	unmarshal := koanf.UnmarshalConf{
		Tag: "koanf",
		DecoderConfig: &mapstructure.DecoderConfig{
			DecodeHook:       decodeHook,
			Metadata:         nil,
			Result:           config,
			WeaklyTypedInput: true,
			TagName:          "koanf",
		},
	}

	if err = ko.UnmarshalWithConf("", config, unmarshal); err != nil {
		return nil, err
	}

	return config, nil
}
