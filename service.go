//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kardianos/service"
	"github.com/yaling888/soulx/common"
)

var (
	soulService, _ = service.New(nil, &service.Config{
		Name:        "Soul",
		DisplayName: "Soul Service",
		Description: "Control the Soul system service.",
	})

	execDir, _ = os.Getwd()
	execPath   = filepath.Join(execDir, fmt.Sprintf("SoulService-%s-%s.exe", runtime.GOOS, runtime.GOARCH))
)

func installService() error {
	_, err := common.ExecCmd(fmt.Sprintf("%s -service install", execPath))
	return err
}

func uninstallService() error {
	_, err := common.ExecCmd(fmt.Sprintf("%s -service uninstall", execPath))
	return err
}

func startService() error {
	_, err := common.ExecCmd(fmt.Sprintf("%s -service start", execPath))
	return err
}

func stopService() error {
	_, err := common.ExecCmd(fmt.Sprintf("%s -service stop", execPath))
	return err
}
