//go:build windows

package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func setPlatformSpecificAttrs(sysAttr *syscall.SysProcAttr) {
	sysAttr.HideWindow = true
	sysAttr.CreationFlags = syscall.CREATE_NEW_PROCESS_GROUP
}

func StartUpdateApp(up *Updater) {
	fh := Resolve[FileHandler](up.DI, "fileHandler")

	exeCutePath := fh.GetExePath()

	cf := Resolve[Config](up.DI, "config")
	up.CurrentExe = up.concatAppName(cf.AppName)
	oldExe := fh.PathJoin(exeCutePath, up.CurrentExe)
	newExe := fh.PathJoin(exeCutePath, "tmp", up.CurrentExe)

	isExist := fh.pathExists(newExe)

	if !isExist {
		_ = up.AppCheckForUpdates()
		_ = up.DownloadUpdates()
	}

	batPath := fh.PathJoin(os.TempDir(), "update_script.bat")
	installDir := filepath.Dir(oldExe)

	batContent := fmt.Sprintf(`
	@echo off
	set "oldExe=%s"
	set "newExe=%s"
	set "installDir=%s"

	:loop
	del /f /q "%%oldExe%%"
	if exist "%%oldExe%%" (
		timeout /t 1 >nul
		goto loop
	)

	move /y "%%newExe%%" "%%oldExe%%"

	rem Change back to the install directory before starting, to avoid config path issues
	cd /d "%%installDir%%"
	start "" "%%oldExe%%"

	rem Self-delete the script
	del "%%~f0"
	`, oldExe, newExe, installDir)

	if err := os.WriteFile(batPath, []byte(batContent), 0644); err != nil {
		return
	}

	cmd := exec.Command("cmd", "/c", batPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP, // create no window
	}

	if err := cmd.Start(); err != nil {
		return
	}
	os.Exit(0)
}
