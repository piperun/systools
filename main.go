package main


import (
        "runtime"
        "os"
        "fmt"
        "bufio"
        "strings"
        "regexp"
        "strconv"
        "log"
        "time"
        "os/exec"
        "golang.org/x/sys/unix"
        "systools/whereis"
)

type System struct {
        OS string
        Host string
        Kernel string
        Uptime map[string]int
        Packages int
        Shell string
        Terminal string
        CPU string
        GPU string
        Memory string
}


type lsb struct {
        ID string
        Release string
        Codename string
        Description string
}

type procvars struct {
        CPU string
        Kernel string

}

type file struct {
        filename string
        exists bool
}



type packages struct {
        num int
        pkgmanager string
}

func getLSBVars() lsb {
        return lsb {
                ID: "DISTRIB_ID",
                Release: "DISTRIB_RELEASE",
                Codename: "DISTRIB_CODENAME",
                Description: "DISTRIB_DESCRIPTION",
        }
}


func procLocations() procvars{
        return procvars {
                Kernel: "/proc/sys/kernel/osrelease",
                CPU: "/proc/cpuinfo",
        }
}



func GetDistro() string{
        const LSB = "/etc/lsb-release"
        const OSREL = "/etc/os-release"
        var distro string
        var LSBvars = getLSBVars()


        switch os := runtime.GOOS; os {
        case "linux":
                distro = getSlice(LSB, LSBvars.Description, "=")
        }
        return distro
}


func GetHost() {

}

func GetCPU() string{
        const model = "model name"
        const cores = "core id"
        const hertz = "cpu MHz"
        const syscpufile = "/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq"
        const cpuinfo = "/proc/cpuinfo"

        var cpufiles [2]string
        cpufiles[0] = "/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq"
        cpufiles[1] = "/proc/cpuinfo"

        re := regexp.MustCompile(`\@\s\d*\.\d*[G-M]Hz`)

        // TODO
        // * remove anything after @, maybe save it later
        // * check if /sys/devices/system/cpu/cpu0 contains cpufreq folder and cpuinfo_min cpuinfo_cur cpuinfo_max
        //  - get max as a default but future config should ask
        // * default to /proc/cpuinfo if nothing found in /sys/devices/system/cpu/cpu0
        /*
        CPU core flow:
        call runtime.NumCPU() -> warn user it might not be accurate
        - if fail (error handling)
        read cpu0 as reffered above
        - if fail (error handling)
        read /proc/cpuinfo

        */

        var cpu string
        var procvars = procLocations()


        cpu = getSlice(procvars.CPU, model, ": ")
        cpu += " @ 2.0GHz"
        match :=  re.FindStringSubmatch(cpu)
        if match[0] != "" {
                cpu = strings.Split(cpu, match[0])[0]
        }

        cpu += "(" + strconv.Itoa(runtime.NumCPU()) + ")"

        for i := 0; i < len(cpufiles); i += 1 {
                if checkFile(cpufiles[i]) {
                        cpu += " @ " + parseCPUFiles(cpufiles[i], hertz)[0:5] + "GHz"
                }
        }



        return cpu
}

func GetKernel() string{
        var uname unix.Utsname
        var version = ""
        unix.Uname(&uname)
        for i := 0; i < len(uname.Release); i += 1 {
                version += string(uname.Release[i])
        }
        return version
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


func GetPackages() {
        var (
                pkgmanager = make(map[string]string)
                f []whereis.File
        )

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
        pkgmanager["gaze"] = "gaze installed" // sorcery
        pkgmanager["alps"] = "alps showinstalled"
        pkgmanager["butch"] = "butch list"
        pkgmanager["bonsai"] = "bonsai list"

        f = getPackagemanager(pkgmanager)
        for i := 0; i < len(f); i += 1 {
                getPackageNum(pkgmanager, f[i])
        }

}


func getPackageNum(pkgmanager map[string]string, file whereis.File) {
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
        fmt.Println(pkgs)

}

func main() {
        var system System
        system.OS = GetDistro()
        system.CPU = GetCPU()
        system.Kernel = GetKernel()
        system.Uptime = GetUptime()
        fmt.Println(system.OS)
        fmt.Println(system.CPU)
        fmt.Println(system.Kernel)
        fmt.Println(system.Uptime)
        GetPackages()
}


func getSlice(filename string, str string, pattern string) string{
        var distro []string
        file, err := os.Open(filename)

        if err != nil {
                panic(err)
        }
       }
        defer file.Close()


        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                if strings.Contains(scanner.Text(), str) {
                        distro = strings.SplitAfter(scanner.Text(), pattern)
                        break
                }
        }
        if len(distro) != 0 {
                return distro[1]
        } else {
                return "NULL - ERROR"
        }

}

func getFile(filename string) []string{
        var content []string

        file, err := os.Open(filename)
        if err != nil {
                log.Fatal(err)
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                content = append(content, scanner.Text())
        }

        if err := scanner.Err(); err != nil {
                log.Fatal(err)
        }
        return content
}

func converttoGHz(hertz string) string{
        var floatHz float64
        var err error

        floatHz, err = strconv.ParseFloat(hertz, 64)
        if strings.Contains(hertz, "."){
                floatHz = floatHz/1000.00
        } else {
                floatHz = floatHz/1000000.00
        }
        if err != nil {
                log.Fatal(err)
        }

        return strconv.FormatFloat(floatHz, 'f', -1, 64)

}

func parseCPUFiles(file string, pattern ...string) string{
        var hertz = ""
        switch f :=  file; f {
        case "/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq":
                hertz = getFile(f)[0]
                return converttoGHz(hertz)
        case "/proc/cpuinfo":
                hertz = getSlice(f, pattern[0], ": ")
                return converttoGHz(hertz)


        default:
                return ""
        }
}

func checkFile(filename string) bool{
        if _, err := os.Stat(filename); os.IsNotExist(err) {
                return false
        }
        return true
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



