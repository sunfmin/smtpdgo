package main

import (
	"strings"
)

func TrimCRLF(s []byte) []byte {
	start, end := 0, len(s)-1
	for start < end {
		if s[end] != '\r' && s[end] != '\n' && s[end] != 0 {
			break
		}
		end -= 1
	}
	return s[start : end+1]
}

func getAddressFromMailFrom(s []byte) []byte {
	l := len(s)
	ab := strings.Index(string(s), "<")
	return s[ab+1 : l-1]
}

func getAddressFromRcptTo(s []byte) []byte {
	l := len(s)
	ab := strings.Index(string(s), "<")
	return s[ab+1 : l-1]
}

func getDomainFromHELO(s []byte) []byte {
	l := len(s)
	return s[5:l]
}
