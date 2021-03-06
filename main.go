package main


import (
	"fmt"
)

type System interface {
	Init()
	GetCPU()
	GetDistro()
	GetGPU()
	GetKernel()
	GetMemory()
	GetModel() Model
	GetOS()
	GetPackages()[]Packages
	GetShell()
	GetTerminal()
	GetUptime()
	GetUserHost()

	// Optional
	GetHDD()
	GetNetwork()
	// Misc
	GetFont()
	GetGTKFont()
	GetDS() // Display Server
}

type SystemForm struct {
	OS string
	UserHost string
	Distro string
	PCModel Model
	Kernel string
	Uptime map[string]int
	Packages []Packages
	Shell Shell
	Terminal string
	CPU string
	GPU string
	Memory map[string]*Storage
}

// These are struct used to store and link methods unique to the OS platform
type BSD struct {
	SystemForm
}

type Linux struct {
	SystemForm
}

type Mac struct {
	SystemForm
}

type Windows struct {
	SystemForm
}




type lsb struct {
	ID string
	Release string
	Codename string
	Description string
}

type file struct {
	filename string
	exists bool
}

type Packages struct {
	num int
	pkgmanager string
}

type Shell struct {
	name string
	version string
	path string
}

type Model struct {
	Name string
	Version string
}

type Storage struct {
	Size float64
	Type string
}

func getLSBVars() lsb {
	return lsb {
		ID: "DISTRIB_ID",
		Release: "DISTRIB_RELEASE",
		Codename: "DISTRIB_CODENAME",
		Description: "DISTRIB_DESCRIPTION",
	}
}


func main() {
	var system System
	system.Distro = GetDistro()
	system.CPU = GetCPU()
	system.Kernel = GetKernel()
	system.Uptime = GetUptime()
	system.Packages = GetPackages()
	system.Shell = GetShell()
	system.Terminal = GetTerminal()
	system.UserHost = GetUserHost()
	system.PCModel = GetModel()
	system.Memory, _ = GetMemory()

	fmt.Println(system.Distro)
	fmt.Println(system.CPU)
	fmt.Println(system.Kernel)
	fmt.Println(system.Uptime)
	fmt.Println(system.Packages)
	fmt.Println(system.Terminal)
	fmt.Println(system.Shell)
	fmt.Println(system.UserHost)
	fmt.Println(system.PCModel)
	fmt.Println(system.Memory["Used"])

}

