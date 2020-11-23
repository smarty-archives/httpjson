package httpjson

import (
	"strconv"
	"strings"
)

func RetrievePrefixedPathInteger(path string, prefixElement string) uint64 {
	value := RetrievePrefixedPathString(path, prefixElement)
	result, _ := strconv.ParseUint(value, 10, 64)
	return result
}

func RetrievePrefixedPathString(path string, prefixElement string) string {
	elements := strings.Split(path, "/")
	for index, element := range elements {
		if element == prefixElement && index < len(elements)-1 {
			return strings.TrimSpace(elements[index+1])
		}
	}
	return ""
}
