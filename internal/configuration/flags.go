package configuration

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/spf13/pflag"
)

func loadFlags(root *koanf.Koanf, delim string, strictMerge bool, flags *pflag.FlagSet) (err error) {
	ko := koanf.NewWithConf(koanf.Conf{
		Delim:       delim,
		StrictMerge: strictMerge,
	})

	if err = ko.Load(posflag.Provider(flags, delim, ko), nil); err != nil {
		return err
	}

	return root.Merge(ko)
}
