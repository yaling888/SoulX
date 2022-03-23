package common

import (
	"os"
	P "path"
	"path/filepath"
)

const Name = "clash"

// Path is used to get the configuration path
var Path = func() *path {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir, _ = os.Getwd()
	}

	homeDir = P.Join(homeDir, ".config", Name)
	return &path{homeDir: homeDir, configFile: "config.yaml"}
}()

type path struct {
	homeDir    string
	configFile string
}

func (p *path) HomeDir() string {
	return p.homeDir
}

func (p *path) ConfigFile() string {
	return filepath.Join(p.homeDir, p.configFile)
}
