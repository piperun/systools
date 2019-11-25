package main

import (
	"regexp"
	"strings"
	"strconv"
	"runtime"
	"log"
)


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


