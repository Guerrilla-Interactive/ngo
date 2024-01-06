package cmd

import (
	"regexp"
	"strings"
)

// Returns the kebabcase version of the title string
func RouteTitleKebabCase(title string) string {
	re := regexp.MustCompile(`\s+`)
	name := strings.ToLower(title)
	// Replace whitespace with -
	name = re.ReplaceAllString(name, "-")
	return name
}
