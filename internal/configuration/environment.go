package configuration

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
)

func loadEnvironment(root *koanf.Koanf, delim string, strictMerge bool) (err error) {
	ko := koanf.NewWithConf(koanf.Conf{
		Delim:       delim,
		StrictMerge: strictMerge,
	})

	if err = ko.Load(env.ProviderWithValue("WEBAUTHN", delim, EnvironmentProviderWithValueCallback), nil); err != nil {
		return err
	}

	return root.Merge(ko)
}
