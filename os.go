package main

import (
	"runtime"
	"golang.org/x/sys/unix"

)




func GetKernel() string{
        var uname unix.Utsname
        var version = ""
        unix.Uname(&uname)
        for i := 0; i < len(uname.Release); i += 1 {
                version += string(uname.Release[i])
        }
        return version
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


