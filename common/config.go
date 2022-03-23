package common

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ExternalController string `yaml:"external-controller" json:"external-controller"`
	Secret             string `yaml:"secret" json:"secret"`
}

func readConfig(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("configuration file %s is empty", path)
	}

	return data, err
}

func Parse() (*Config, error) {
	return ParseWithPath(Path.ConfigFile())
}

func ParseWithPath(path string) (*Config, error) {
	buf, err := readConfig(path)
	if err != nil {
		return nil, err
	}

	return ParseConfig(buf)
}

func ParseConfig(buf []byte) (*Config, error) {
	config := &Config{}

	if err := yaml.Unmarshal(buf, config); err != nil {
		return nil, err
	}

	return config, nil
}
