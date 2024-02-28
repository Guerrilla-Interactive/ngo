package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	errMissingLeadingSlashInRouteName = errors.New("route name must contain leading /")
	errTrailingSlashInRouteName       = errors.New("route name must not contain a trailing /")
	errMultipleSlashInRouteName       = errors.New("route must not contain multiple //")
	errWhitespaceInRouteName          = errors.New("route must not contain any whitespace")
	errInvalidFolderName              = errors.New("folder name is invalid")
	errCannotAddChildrenToStaticRoute = errors.New("cannot add children to static route")
)

var (
	FillerRouteNameRegex = regexp.MustCompile(`^\(([[[:alpha:]]|-)+[[:alpha:]]+\)$`)
	StaticRouteNameRegex = regexp.MustCompile(`^([[[:alpha:]]|-)+[[:alpha:]]+$`)

	DynamicRouteNameRegex                 = regexp.MustCompile(`^\[slug\]$`)
	DynamicRouteCatchAllNameRegex         = regexp.MustCompile(`^\[\.\.\.slug\]$`)
	DynamicRouteOptionalCatchAllNameRegex = regexp.MustCompile(`^\[\[\.\.\.slug\]\]$`)

	// Note the general versions are different than the previous in that general don't enfore
	// pattern to match from beginning to string to the end
	GeneralDynamicRouteNameRegex                 = regexp.MustCompile(`\[[[:alnum:]]+\]`)
	GeneralDynamicRouteCatchAllNameRegex         = regexp.MustCompile(`\[\.\.\.[[:alnum:]]+\]`)
	GeneralDynamicRouteOptionalCatchAllNameRegex = regexp.MustCompile(`\[\[\.\.\.[[:alnum:]]+\]\]`)
	IndexRouteEnding                             = "(index)" //
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

// Returns whether a given folder name is valid. This function may also be used to test
// whether a given route part is valid
// Current limitation is that only alphabet and dashes are allowed in route parts or folders.
// This must be changed as necessary
func IsValidFolderName(name string) bool {
	nameByte := []byte(name)
	startsWithAlphabet := regexp.MustCompile(`^[[:alpha:]]`).Match(nameByte)
	endsWithAlphabet := regexp.MustCompile(`[[:alpha:]]$`).Match(nameByte)
	alphabetAndDashes := regexp.MustCompile(`([[:alpha:]]|-)+`).Match(nameByte)
	return startsWithAlphabet && endsWithAlphabet && alphabetAndDashes
}

// Checks if the route name is valid.
// If invalid, the returned error message contains details
// The function IsRouteNameValid does the same job but returns a
// boolean indicating if the route name is valid instead of an
// error object.
func RouteNameValid(name string) error {
	// Root route
	if len(name) == 0 {
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
	// static index route
	if name == fmt.Sprintf("/%v", IndexRouteEnding) {
		return nil
	}
	// dynamic index route
	if name == "/[slug]" || name == "/[...slug]" || name == "[[...slug]]" {
		return nil
	}
	// Filler route. Note this is different from a filler folder whose name
	// starts and ends in parenthesis.
	if IsValidFolderName(name) {
		return nil
	}

	validSuffixes := []string{
		fmt.Sprintf("/%v", IndexRouteEnding),
		"/[slug]",
		"/[...slug]",
		"/[[...slug]]",
	}
	for _, suffix := range validSuffixes {
		if strings.HasSuffix(name, suffix) {
			parentRoute := strings.TrimSuffix(name, suffix)
			kind, err := RouteTypeFromRouteName(parentRoute)
			if err != nil {
				return err
			}
			if kind == StaticRoute {
				return errCannotAddChildrenToStaticRoute
			}
			return nil
		}
	}

	routeParts := strings.Split(name, "/")
	lastPart := routeParts[len(routeParts)-1]
	if IsValidFolderName(lastPart) {
		parentRoute := strings.Join(routeParts[:len(routeParts)-1], "/")
		kind, err := RouteTypeFromRouteName(parentRoute)
		if err != nil {
			return err
		}
		if kind == StaticRoute {
			return errCannotAddChildrenToStaticRoute
		}
		return nil
	} else {
		return errInvalidFolderName
	}
}

// Returns a boolean indicating if a given route name is valid
// The function AssertRouteName does the same job but returns an
// non-nil error when the given name is an invalid route name
func IsRouteNameValid(name string) bool {
	return RouteNameValid(name) == nil
}

// Return the route type based on the given candiate string
// also return non nil error when the input is invalid (i.e. no route type
// exists for the given name)
// Preconditions:
// name is a valid route name
func RouteTypeFromRouteName(name string) (RouteType, error) {
	err := RouteNameValid(name)
	if err != nil {
		return *new(RouteType), err
	}
	if name == "" {
		return RootRoute, nil
	}
	if strings.HasSuffix(name, fmt.Sprintf("/%v", IndexRouteEnding)) {
		return StaticRoute, nil
	} else if strings.HasSuffix(name, "/[slug]") {
		return DynamicRoute, nil
	} else if strings.HasSuffix(name, "/[...slug]") {
		return DynamicCatchAllRoute, nil
	} else if strings.HasSuffix(name, "/[[...slug]]") {
		return DynamicCatchAllRoute, nil
	} else {
		return FillerRoute, nil
	}
}

func GetSchemaExportName(routeName string, routeType RouteType) (string, error) {
	nameVar := GetRouteTemplateVariable(routeName)
	switch routeType {
	case StaticRoute:
		return fmt.Sprintf("%vIndexSchema", nameVar.CamelCaseComponentName), nil
	case DynamicRoute:
		return fmt.Sprintf("%vSlugSchema", nameVar.CamelCaseComponentName), nil
	case DynamicCatchAllRoute:
		return fmt.Sprintf("%vSlugCASchema", nameVar.CamelCaseComponentName), nil
	case DynamicCatchAllOptionalRoute:
		return fmt.Sprintf("%vSlugCAOSchema", nameVar.CamelCaseComponentName), nil
	default:
		return "", errors.New("not implemented for given route")
	}
}
