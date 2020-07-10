package utils

import "strings"

func StringInSlice(needle string, haystack []string) bool {
	for _, hay := range haystack {
		if hay == needle {
			return true
		}
	}
	return false
}

func SuffixInSlice(needle string, haystack []string) bool {
	for _, h := range haystack {
		if strings.HasSuffix(needle, h) {
			return true
		}
	}
	return false
}

func StringContainedInSlice(needle string, haystack []string) bool {
	for _, h := range haystack {
		if strings.Contains(needle, h) {
			return true
		}
	}

	return false
}
