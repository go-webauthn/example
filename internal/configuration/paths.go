package configuration

import (
	"os"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

func loadPaths(root *koanf.Koanf, paths []string, delim string, strictMerge bool) (err error) {
	for _, path := range paths {
		err = loadPath(root, path, delim, strictMerge)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadPath(root *koanf.Koanf, path string, delim string, strictMerge bool) (err error) {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return loadPathDir(root, path, delim, strictMerge)
	}

	return loadPathFile(root, path, delim, strictMerge)
}

func loadPathFile(root *koanf.Koanf, path string, delim string, strictMerge bool) (err error) {
	parser := loadPathFileParser(path)

	return loadPathFileWithParser(root, parser, path, delim, strictMerge)
}

func loadPathFileParser(path string) (parser koanf.Parser) {
	ext := filepath.Ext(path)

	switch ext {
	case ".yml", ".yaml":
		return yaml.Parser()
	case ".tml", ".toml":
		return toml.Parser()
	case ".json":
		return json.Parser()
	default:
		return nil
	}
}

func loadPathFileWithParser(root *koanf.Koanf, parser koanf.Parser, path string, delim string, strictMerge bool) (err error) {
	ko := koanf.NewWithConf(koanf.Conf{
		Delim:       delim,
		StrictMerge: strictMerge,
	})

	if err = ko.Load(file.Provider(path), parser); err != nil {
		return err
	}

	return root.Merge(ko)
}

func loadPathDir(root *koanf.Koanf, path string, delim string, strictMerge bool) (err error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			err := loadPathDir(root, filepath.Join(path, entry.Name()), delim, strictMerge)
			if err != nil {
				return err
			}
		} else {
			parser := loadPathFileParser(filepath.Join(path, entry.Name()))
			if parser == nil {
				continue
			}

			err := loadPathFile(root, filepath.Join(path, entry.Name()), delim, strictMerge)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
