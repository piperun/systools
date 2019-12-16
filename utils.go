package main


import (
	"strings"
	"os"
	"log"
	"bufio"
	"math"
)


func RemovePrefix(str string, prefixes ...string) string{
	for _, prefix := range prefixes {
		str = strings.TrimPrefix(str, prefix)

	}
	return str
}

func RemoveSuffix(str string, suffixes ...string) string{
	for _, suffix := range suffixes {
		str = strings.TrimSuffix(str, suffix)

	}
	return str
}

func (s *Storage) Conversion(bytes uint64) {
	var metrics = []string{
		"Bytes",
		"KiB",
		"MiB",
		"GiB",
		"TiB",
		"PiB",
		"EiB",
		"ZiB",
		"YiB",
	}

	exponent := math.Round(float64(DigitCount(bytes))/4)
	s.Size = float64(bytes) / math.Pow(1024, exponent)
	s.Type = metrics[int(exponent)]
}

func DigitCount(i uint64) uint64{
	var count uint64 = 0
	for i != 0 {
		i /= 10
		count += 1
	}
	return count
}

func GetSlice(filename string, str string, pattern string) string{
	var slice []string
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), str) {
			slice = strings.SplitAfter(scanner.Text(), pattern)
			break
		}
	}
	if len(slice) != 0 {
		return slice[1]
	} else {
		return "NULL - ERROR"
	}
}
