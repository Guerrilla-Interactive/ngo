package main

import (
	"regexp"
	"strings"
)

// Return true if s is y, Y or Yes (in any case) with any white-space leading
// or subsequent white-space
func gotYes(s string) bool {
	s = strings.ToUpper(s)
	return regexp.MustCompile(`^\s*(?:Y|YES)\s*$`).MatchString(s)
}
