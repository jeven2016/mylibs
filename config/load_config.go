package config

import (
	"github.com/jeven2016/mylibs/internal"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"log"
)

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var k = koanf.New(".")

// LoadConfig loads the configuration files
func LoadConfig(internalCfg []byte, config Config, extraConfigFilePath *string, defaultCfgPaths []string) error {

	//load internal config
	if internalCfg != nil {
		if err := k.Load(rawbytes.Provider(internalCfg), yaml.Parser()); err != nil {
			return err
		}
	}

	cfgPaths := defaultCfgPaths
	if cfgPaths == nil {
		cfgPaths = []string{}
	}
	if extraConfigFilePath != nil {
		cfgPaths = append(cfgPaths, *extraConfigFilePath)
	}

	// load external configs
	for _, f := range cfgPaths {
		if exists, err := internal.IsFileExists(f); err != nil {
			log.Printf(f + " not found and ignored")
			continue
		} else if exists {
			if err = k.Load(file.Provider(f), yaml.Parser()); err != nil {
				return err
			}
		}
	}

	if err := k.Unmarshal("", config); err != nil {
		return err
	}

	return nil
}
