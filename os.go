package main

import (
	"runtime"
	"golang.org/x/sys/unix"
)


func GetOS() string {
	return runtime.GOOS
}

func GetKernel() string {
	var (
		uname unix.Utsname
		version string
	)

	unix.Uname(&uname)
	for i := 0; i < len(uname.Release); i += 1 {
		version += string(uname.Release[i])
	}
	return version
}

func (linux *Linux) GetDistro() {
	const(
		osrelvar = "PRETTY_NAME"
		lsbvar = "DISTRIB_DESCRIPTION"
		pattern = "="
	)
	var (
		distro string
		releaseFiles = []string {
			"/etc/os-release",
			"/etc/lsb-release",
			"/etc/bedrock-release",
		}
	)
	for i := 0; i < len(releaseFiles); i += 1 {
		if CheckFile(releaseFiles[i]) {
			switch releaseFiles[i] {
			case "/etc/os-release":
				distro = GetSlice(releaseFiles[i], osrelvar, "=")
			case "/etc/lsb-release":
				distro = GetSlice(releaseFiles[i], lsbvar, "=")
			}
			break
		}
	}
	linux.Distro = distro

}


