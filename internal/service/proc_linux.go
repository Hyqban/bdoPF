//go:build linux

package service

import (
	"os"
	"os/exec"
	"syscall"
)

func setPlatformSpecificAttrs(sysAttr *syscall.SysProcAttr) {
	sysAttr.Setpgid = true
}

func StartUpdateApp(up *Updater) {
	fh := Resolve[FileHandler](up.DI, "fileHandler")

	exxCutePath := fh.GetExePath()

	cf := Resolve[Config](up.DI, "config")
	up.CurrentExe = up.concatAppName(cf.AppName)
	oldExe := fh.PathJoin(exxCutePath, up.CurrentExe)
	newExe := fh.PathJoin(exxCutePath, "tmp", up.CurrentExe)

	isExist := fh.pathExists(newExe)

	if !isExist {
		_ = up.AppCheckForUpdates()
		_ = up.DownloadUpdates()
	}

	if err := os.Chmod(newExe, 0755); err != nil {
		return
	}

	if err := os.Rename(newExe, oldExe); err != nil {
		return
	}

	args := os.Args
	env := os.Environ()

	err := syscall.Exec(oldExe, args, env)
	if err != nil {
		cmd := exec.Command(oldExe, args[1:]...)
		cmd.Start()
		os.Exit(0)
	}
}
