package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func YamlDeCode(b []byte) (Config, error) {
	var c Config
	err := yaml.Unmarshal(b, &c)
	if err != nil {
		return c, fmt.Errorf("YamlDeCode: %w", err)
	}
	return c, nil
}
