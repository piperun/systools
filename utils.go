package main


import (
	"strings"
)


func RemovePrefix(str string, prefixes ...string) string{
	for _, prefix := range prefixes {
		str = strings.TrimPrefix(str, prefix)

	}
	return str
}
