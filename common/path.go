package common

import (
	"os"
	P "path"
	"path/filepath"
)

const Name = "clash"

var ExecutableDir = getExecutableDir()

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

func (p *path) MMDB() string {
	return P.Join(p.homeDir, "Country.mmdb")
}

func (p *path) GeoSite() string {
	return P.Join(p.homeDir, "geosite.dat")
}

func (p *path) ExternalUI() string {
	return P.Join(p.homeDir, "dashboard")
}

func (p *path) Resolve(path string) string {
	if !filepath.IsAbs(path) {
		return filepath.Join(p.homeDir, path)
	}

	return path
}

func getExecutableDir() string {
	fullExecPath, _ := os.Executable()
	dir, _ := filepath.Split(fullExecPath)
	return dir
}

func ResolveExecutableResourcesDir(path string) string {
	if !filepath.IsAbs(path) {
		return filepath.Join(ExecutableDir, "resources", path)
	}

	return path
}
