package main

import (
	"regexp"
	"strings"
	"strconv"
	"runtime"
	"log"
)


func GetCPU() string{
	const (
		model = "model name"
		cores = "core id"
		hertz = "cpu MHz"
		syscpufile = "/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq"
		cpuinfo = "/proc/cpuinfo"
		minlen = 5
	)
	var (
		cpufiles [2]string
		cpu string
		procvars = procLocations()
	)

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

	cpu = getSlice(procvars.CPU, model, ": ")
	match :=  re.FindStringSubmatch(cpu)

	if match != nil {
		cpu = strings.Split(cpu, match[0])[0]
	}

	cpu += "(" + strconv.Itoa(runtime.NumCPU()) + ")"

	for i := 0; i < len(cpufiles); i += 1 {
		if checkFile(cpufiles[i]) {
			tmp := parseCPUFiles(cpufiles[i], hertz)
			if len(tmp) > minlen {
				cpu += " @ " + tmp[0:minlen] + "GHz"
			} else if len(tmp) < minlen  {
				cpu += " @ " + tmp + "GHz"

			}

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


