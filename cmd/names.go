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
	FillerRouteNameRegex = regexp.MustCompile(`^\([[[:alnum:]]|-]+\)$`)
	StaticRouteNameRegex = regexp.MustCompile(`^[[[:alpha:]]|-]+$`)

	DynamicRouteNameRegex                 = regexp.MustCompile(`^\[slug\]$`)
	DynamicRouteCatchAllNameRegex         = regexp.MustCompile(`^\[\.\.\.slug\]$`)
	DynamicRouteOptionalCatchAllNameRegex = regexp.MustCompile(`^\[\[\.\.\.slug\]\]$`)

	// Note the general versions are different than the previous in that general don't enfore
	// pattern to match from beginning to string to the end
	GeneralDynamicRouteNameRegex                 = regexp.MustCompile(`\[slug\]`)
	GeneralDynamicRouteCatchAllNameRegex         = regexp.MustCompile(`\[\.\.\.slug\]`)
	GeneralDynamicRouteOptionalCatchAllNameRegex = regexp.MustCompile(`\[\[\.\.\.slug\]\]`)
)

// Returns the kebabcase version of the title string
func RouteTitleKebabCase(title string) string {
	name := strings.ToLower(title)
	re := regexp.MustCompile(`\s+`)

	// Remote leading whitespace
	leadingWS := regexp.MustCompile(`^\s+`)
	name = leadingWS.ReplaceAllString(name, "")

	// Remote trailing whitespace
	trailingWS := regexp.MustCompile(`\s+$`)
	name = trailingWS.ReplaceAllString(name, "")

	// Replace other whitespace with -
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
	if name == "/" && kind == StaticRoute {
		return nil
	}
	if !strings.HasPrefix(name, "/") {
		return errMissingLeadingSlashInRouteName
	}
	if strings.HasSuffix(name, "/") && name != "/" {
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

	// There should be no filler routes (except at last?)
	// The input to "ngo add route" has to be as similar as possible
	// of the "ngo ls" command
	for i := 0; i < len(routeParts)-1; i++ {
		if routeParts[i] == "" {
			continue
		}
		if IsValidFillerRouteName(routeParts[i]) {
			return fmt.Errorf("route name cannot contain filler route got %q", routeParts[i])
		} else if !IsValidStaticRouteName(routeParts[i]) && !IsValidDynamicRouteName(routeParts[i]) {
			return fmt.Errorf("invalid route name %q", routeParts[i])
		}
	}

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
			return errors.New("dynamic route must end with [slug] or [...slug] or [[...slug]]")
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

// Return the route type based on the given candiate string
// also return non nil error when the input is invalid (i.e. no route type
// exists for the given name)
func praseRouteType(candidate string) (RouteType, error) {
	candidateLower := strings.ToLower(candidate)
	switch candidateLower {
	case "static":
		return StaticRoute, nil
	case "dynamic":
		return DynamicRoute, nil
	case "filler":
		return FillerRoute, nil
	}
	return StaticRoute, fmt.Errorf("invalid route type. valid types are static, dynamic, filler got %q", candidate)
}

// Return the route name based on the given candiate string
// also return non nil error for when the input is invalid
func praseRouteName(candidate string) (string, error) {
	return RouteTitleKebabCase(candidate), nil
}

func GetSchemaExportName(routeName string, routeType RouteType) (string, error) {
	nameVar := GetRouteTemplateVariable(routeName)
	switch routeType {
	case StaticRoute:
		return fmt.Sprintf("%vIndexSchema", nameVar.CamelCaseComponentName), nil
	case DynamicRoute:
		return fmt.Sprintf("%vSlugSchema", nameVar.CamelCaseComponentName), nil
	default:
		return "", errors.New("not implemented for given route")
	}
}
