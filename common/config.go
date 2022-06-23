package common

import (
	"fmt"
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ExternalController string `yaml:"external-controller" json:"external-controller"`
	ExternalUI         string `yaml:"external-ui" json:"external-ui"`
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

	if host, port, err := net.SplitHostPort(config.ExternalController); err == nil {
		if ip := net.ParseIP(host); ip != nil && ip.IsUnspecified() {
			config.ExternalController = net.JoinHostPort("127.0.0.1", port)
		}
	}

	return config, nil
}

func InitResourcesIfNotExist() {
	dir := Path.HomeDir()

	if _, err := os.Stat(Path.MMDB()); os.IsNotExist(err) {
		_, _ = ExtractZipFile(ResolveExecutableResourcesDir("Country.mmdb.zip"), dir)
	}

	if _, err := os.Stat(Path.GeoSite()); os.IsNotExist(err) {
		_, _ = ExtractZipFile(ResolveExecutableResourcesDir("geosite.dat.zip"), dir)
	}

	if _, err := os.Stat(Path.ExternalUI()); os.IsNotExist(err) {
		_, _ = ExtractZipFile(ResolveExecutableResourcesDir("dashboard.zip"), dir)
	}
}
