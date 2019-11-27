package main

import (
	"strconv"
	"strings"
	"os"
	"fmt"
	"os/exec"
	"time"
	"golang.org/x/sys/unix"
	"github.com/piperun/systools/whereis"
)



func GetHost() {

}


func GetTerminal() {
}

func GetShell() {
	var (
		shell string
		cmd *exec.Cmd
	)

	shell = os.Getenv("SHELL")
	cmd = exec.Command(shell, "--version")
	temp, _ := cmd.Output()
	output := string(temp)

	fmt.Println(output)



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

func GetPackage() {

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
