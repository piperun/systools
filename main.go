package main


import (
        "os"
        "fmt"
        "bufio"
        "strings"
        "log"
)

type System struct {
        OS string
        Host string
        Kernel string
        Uptime map[string]int
        Packages []packages
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


func main() {
        var system System
        system.OS = GetDistro()
        system.CPU = GetCPU()
        system.Kernel = GetKernel()
        system.Uptime = GetUptime()
	system.Packages = GetPackages()
        fmt.Println(system.OS)
        fmt.Println(system.CPU)
        fmt.Println(system.Kernel)
        fmt.Println(system.Uptime)
	fmt.Println(system.Packages)
	GetShell()
	GetTerminal()
}


func getSlice(filename string, str string, pattern string) string{
        var distro []string
        file, err := os.Open(filename)

        if err != nil {
                panic(err)
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

func checkFile(filename string) bool{
        if _, err := os.Stat(filename); os.IsNotExist(err) {
                return false
        }
        return true
}


