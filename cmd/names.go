package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	errZeroLengthRouteName            = errors.New("route name must be contain as least one character")
	errMissingLeadingSlashInRouteName = errors.New("route name must contain leading /")
	errTrailingSlashInRouteName       = errors.New("route name must not contain a trailing /")
	errMultipleSlashInRouteName       = errors.New("route must not contain multiple //")
	errWhitespaceInRouteName          = errors.New("route must not contain any whitespace")
)

var (
	FillerRouteNameRegex                  = regexp.MustCompile(`^\([[:alnum:]]+\)$`)
	StaticRouteNameRegex                  = regexp.MustCompile(`^[[:alpha:]]+$`)
	DynamicRouteNameRegex                 = regexp.MustCompile(`^\[[[:alnum:]]+\]$`)
	DynamicRouteCatchAllNameRegex         = regexp.MustCompile(`^\[\.\.\.[[:alnum:]]+\]$`)
	DynamicRouteOptionalCatchAllNameRegex = regexp.MustCompile(`^\[\[\.\.\.[[:alnum:]]+\]\]$`)
)

// Returns the kebabcase version of the title string
func RouteTitleKebabCase(title string) string {
	re := regexp.MustCompile(`\s+`)
	name := strings.ToLower(title)
	// Replace whitespace with -
	name = re.ReplaceAllString(name, "-")
	return name
}

// Checks if the route name is valid.
// If invalid, the returned error message contains details
func AssertRouteNameValid(kind RouteType, name string) error {
	if len(name) == 0 {
		return errZeroLengthRouteName
	}
	// Valid route name
	// Note that we check for this before checking leading and trailing
	// because we require a leading slash but not trailing so its easy
	// to write code that considers '/' as invalid route name but it's
	// a valid static root route.
	if name == "/" {
		return nil
	}
	if !strings.HasPrefix(name, "/") {
		return errMissingLeadingSlashInRouteName
	}
	if strings.HasSuffix(name, "/") {
		return errTrailingSlashInRouteName
	}
	if strings.Contains(name, "//") {
		return errMultipleSlashInRouteName
	}
	if regexp.MustCompile(`\s`).MatchString(name) {
		return errWhitespaceInRouteName
	}
	routeParts := strings.Split(name, "/")
	lastPart := routeParts[len(routeParts)-1]
	// TODO
	// What should we say about filler routes in the route name?
	// Should they exist?
	switch kind {
	case StaticRoute:
		if !IsValidStaticRouteName(lastPart) {
			return fmt.Errorf("%s is not a valid static route name", lastPart)
		}
	case DynamicRoute:
		// Dynamic Route [slug]
		if !IsValidDynamicRouteName(lastPart) {
			return errors.New("dynamic route must end with [alphanum] or [...alphanum] or [[...alphanum]] where alphanum can be replaced with word of your choice")
		}
	case FillerRoute:
		if !IsValidFillerRouteName(lastPart) {
			return errors.New("filler route name must end with (xxx) where xxx may be any alpha-numeric string")
		}
	}
	return nil
}

func IsValidStaticRouteName(candidate string) bool {
	return StaticRouteNameRegex.Match([]byte(candidate))
}

func IsValidFillerRouteName(candidate string) bool {
	return FillerRouteNameRegex.Match([]byte(candidate))
}

func IsValidDynamicRouteName(candidate string) bool {
	return DynamicRouteNameRegex.Match([]byte(candidate)) ||
		DynamicRouteCatchAllNameRegex.Match([]byte(candidate)) ||
		DynamicRouteOptionalCatchAllNameRegex.Match([]byte(candidate))
}

func GetDynamicRouteType(candidate string) (DynamicRouteType, error) {
	if DynamicRouteNameRegex.Match([]byte(candidate)) {
		return DynamicRoutePrimary, nil
	} else if DynamicRouteCatchAllNameRegex.Match([]byte(candidate)) {
		return DynamicRouteCatchAll, nil
	} else if DynamicRouteOptionalCatchAllNameRegex.Match([]byte(candidate)) {
		return DynamicRouteOptionalCatchAll, nil
	}
	var bogus DynamicRouteType
	return bogus, errors.New("not a valid dynamic route type name")
}
