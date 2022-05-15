//go:build windows

package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var homeDir, _ = GetHomeDir()

func LogOutFile() string {
	p := filepath.Join(homeDir, "log")
	_, err := os.Stat(p)
	if err != nil {
		_ = os.MkdirAll(p, os.ModePerm)
	}
	return filepath.Join(homeDir, "log", "soul_out.log")
}

func LogErrFile() string {
	return filepath.Join(homeDir, "log", "soul_err.log")
}

func ExecutableDir() string {
	fullExecPath, _ := os.Executable()
	dir, _ := filepath.Split(fullExecPath)
	return dir
}

func CoreExecPath() string {
	var (
		wd     = ExecutableDir()
		corDir = filepath.Join(wd, "core")
		execs  []string
	)
	_ = filepath.Walk(corDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".exe" {
			return nil
		}
		execs = append(execs, path)
		return nil
	})

	if len(execs) > 0 {
		return execs[0]
	}

	return filepath.Join(wd, "core", fmt.Sprintf("soul-%s-%s.exe", runtime.GOOS, runtime.GOARCH))
}

func GetHomeDir() (string, error) {
	wd := ExecutableDir()
	p := filepath.Join(wd, "path.conf")

	f, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		_ = logger.Warningf("Failed to open std err %q: %v", p, err)
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	content, err := io.ReadAll(f)

	pa := string(content)
	if strings.TrimSpace(pa) == "" {
		pa, _ = os.UserHomeDir()
		pa = filepath.Join(pa, ".config", "clash")
		_, _ = io.WriteString(f, pa)
	}

	return pa, nil
}
