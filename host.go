package main

//#include <unistd.h>
import "C"

import (
	"strconv"
	"strings"
	//"fmt"
	"os"
	"os/user"
	"os/exec"
	"time"
	"golang.org/x/sys/unix"
	"github.com/piperun/systools/whereis"

)


func GetUser() *user.User {
	user, _ := user.Current()
	return user
}

func GetHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func GetModel() Model{
	var model Model
	model.grabModel()
	model.removeDummyinfo()
	return model

}

func GetUserHost() string{
	return GetUser().Username + "@" + GetHostname()
}

func GetTerminal() string {
	return C.GoString(C.ttyname(C.int(unix.Stdin)))
}

func GetShell() Shell{
	var (
		usershell Shell
	)
	usershell.grabShell()
	usershell.grabShellVersion()

	return usershell
}



func GetUptime() map[string]int{
	var (
		sysinfo unix.Sysinfo_t
		date time.Duration
		uptime map[string]int = make(map[string]int)
	)

	unix.Sysinfo(&sysinfo)
	date, _ = time.ParseDuration(strconv.FormatInt(sysinfo.Uptime, 10)+"s")
	uptime["days"] = int(date.Hours()) / 24
	uptime["hours"] = int(date.Hours()) % 24
	uptime["minutes"] = int(date.Minutes()) % 60
	uptime["seconds"] = int(date.Seconds()) % 60
	return uptime
}

func GetPackages() []packages{
	var (
		pkgmanager = make(map[string]string)
		f []whereis.File
		pkgs []packages
	)


	// Maybe switch to slice since map is a bit redundant
	pkgmanager["kiss"] = "kiss l"
	pkgmanager["pacman"] = "pacman -Qq --color never"
	pkgmanager["dpkg"] = "dpkg-query -W --showformat=."
	pkgmanager["rpm"] = "rpm -qa"
	pkgmanager["xbps-query"] =  "xbps-query -l"
	pkgmanager["apk"] = "apk info"
	pkgmanager["opkg"] = "opkg list-installed"
	pkgmanager["pacman-g2"] = "pacman-g2 -Q"
	pkgmanager["lvu"] = "lvu installed"
	pkgmanager["tce-status"] = "tce-status -i"
	pkgmanager["pkg_info"] = "pkg_info"
	pkgmanager["tazpkg"] = "tazpkg list" // subtract -6 from packages according to neofetch's src
	pkgmanager["sorcery"] = "gaze installed" // sorcery
	pkgmanager["alps"] = "alps showinstalled"
	pkgmanager["butch"] = "butch list"
	pkgmanager["bonsai"] = "bonsai list"

	f = getPackagemanager(pkgmanager)
	for i := 0; i < len(f); i += 1 {
		pkgs = append(pkgs, getPackageNum(pkgmanager, f[i]))
	}

	return pkgs

}


func getPackageNum(pkgmanager map[string]string, file whereis.File) packages{
	var (
		pkgs packages
		cmd *exec.Cmd
		output []byte
		cmdline []string
	)

	cmdline = strings.Split(pkgmanager[file.Name], " ")
	cmd = exec.Command(cmdline[0], cmdline[1:]...)
	output, err := cmd.Output()
	if err == nil {
		pkgs.num = len(output)
		pkgs.pkgmanager = file.Name
	}
	return pkgs

}

func getPackagemanager(pkgmanager map[string]string) []whereis.File{
	var (
		paths whereis.Paths
		f whereis.File
		err error
		pkgmanagers []whereis.File
	)

	paths.SetPATHlist()

	for k, _ := range pkgmanager {
		f, err = whereis.Lookup(paths.Pathlist, k)
		if err == nil {
			pkgmanagers = append(pkgmanagers, f)
		}
	}
	return pkgmanagers
}

// Version 1

func (s *Shell) grabShell() {
	var (
		shellexec, shellpath string
		prefixes, shellenvs []string

	)
	prefixes = []string{"/bin/", "/snap/bin/"}
	shellenvs = []string {"SHELL", "_"}


	for _, env := range shellenvs {
		shellpath = os.Getenv(env)
		if shellpath != "" {
			shellexec = RemovePrefix(shellpath, prefixes...)
			break
		}
	}

	s.name = shellexec
	s.path = shellpath


	/*
	if len(shell) == 0
	getshellfrompasswd(getuser()) -> return shell
	*/


}

func (s *Shell) grabShellVersion() {

	var (
		cmd *exec.Cmd
		shellcmd map[string]string = make(map[string]string)
	)

	shellcmd["bash"] = "printf %s $BASH_VERSION"
	shellcmd["ksh"] = "printf %s \"$KSH_VERSION\""
	shellcmd["tcsh"] = "$tcsh"
	shellcmd["base"] = "-c"


	cmd = exec.Command(s.name, shellcmd["base"], shellcmd[s.name])
	temp, _ := cmd.Output()
	s.version = string(temp)

}

func (m *Model) grabModel() {
	var (
		modelfiles []string
	)

	modelfiles = []string{
		"/sys/devices/virtual/dmi/id/product_name",
		"/sys/devices/virtual/dmi/id/product_version",
		"/sys/firmware/devicetree/base/model",
		"/tmp/sysinfo/model",
	}

	for i := 0; i < len(modelfiles); i += 1 {
		if !checkFile(modelfiles[i]) {
			modelfiles[len(modelfiles)-1], modelfiles[i] = modelfiles[i], modelfiles[len(modelfiles)-1]
			modelfiles = modelfiles[:len(modelfiles)-1]

		}
		switch modelfiles[i] {
		case "/sys/devices/virtual/dmi/id/product_name":
			m.Name = getFile(modelfiles[i])[0]
		case "/sys/devices/virtual/dmi/id/product_version":
			m.Version = getFile(modelfiles[i])[0]
		case "/sys/firmware/devicetree/base/model":
			// This is for ARM and I don't have an ARM device.
		case "/tmp/sysinfo/model":
		default:
			m.Name = "NO MODEL NAME FOUND"
			m.Version = "NO MODEL VERSION FOUND"

		}
	}
}

// Used to remove any "dummy"/vm model info
func (m *Model) removeDummyinfo() {
	var (
		dummyinfo []string
	)

	dummyinfo = []string{
		"To be filled by O.E.M.",
		"To Be Filled*",
		"OEM*",
		"Not Applicable",
		"System Product Name",
		"System Version",
		"Undefined",
		"Default string",
		"Not Specified",
		"Type1ProductConfigId",
		"INVALID",
		"�",
	}

	if strings.Contains(m.Name, "Standard PC") {
		m.Name = strings.Replace(m.Name, "Standard PC", "KVM/QEMU", -1)
	} else if strings.Contains(m.Name, "OpenBSD") {
		m.Name = strings.Replace(m.Name, "Standard PC", "vmm", -1)
	}


	for i := 0; i < len(dummyinfo); i += 1 {
		if strings.Contains(m.Name, dummyinfo[i]) {
			strings.Replace(m.Name, dummyinfo[i], "", -1)
		}

	}


}
