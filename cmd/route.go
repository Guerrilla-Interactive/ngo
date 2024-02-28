package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	errPagePathInsufficientLength = errors.New("given path to page.tsx has to contain as least one parent directory")
	errPagePathInvalidName        = errors.New("pagename can be page.tsx or page.jsx")
)

// Checks if the given name is valid route terminal name
// Returns true iff candidate is 'page.tsx' or 'page.jsx'
func IsValidTerminalPageRouteName(candidate string) bool {
	return candidate == "page.tsx" || candidate == "page.jsx"
}

// Preconditions,
// name is a valid folder name for Dynamic or a Static route
func FolderNameToRouteType(name string) RouteType {
	// If foldername is the path to the folder, extract
	// just the last part so we get the folder name instead
	routeParts := strings.Split(name, string(os.PathSeparator))
	lastPart := routeParts[len(routeParts)-1]
	if strings.HasPrefix(lastPart, "[[") && strings.HasSuffix(lastPart, "]]") {
		return DynamicCatchAllOptionalRoute
	} else if strings.HasPrefix(lastPart, "[...") && strings.HasSuffix(lastPart, "]") {
		return DynamicCatchAllRoute
	} else if strings.HasPrefix(lastPart, "[") && strings.HasSuffix(lastPart, "]") {
		return DynamicRoute
	} else {
		return StaticRoute
	}
}

// Preconditions:
// path is a valid full path to page.tsx or page.jsx on
// a nextJS project with app directory
func RouteTypeByPageTSXPath(path string) (RouteType, error) {
	var kind RouteType
	routeParts := strings.Split(path, string(os.PathSeparator))
	lastPart := routeParts[len(routeParts)-1]
	if len(routeParts) < 2 {
		return kind, errPagePathInsufficientLength
	}
	if !IsValidTerminalPageRouteName(lastPart) {
		return kind, errPagePathInvalidName
	}
	// Now we traverse up from the lowest child ignoring all filler directories
	for i := len(routeParts) - 2; i >= 0; i-- {
		name := routeParts[i]
		if isValidFillerDirectory(name) {
			continue
		}
		kind := FolderNameToRouteType(name)
		return kind, nil
	}
	return StaticRoute, nil
}

func isValidFillerDirectory(candidate string) bool {
	return strings.HasPrefix(candidate, "(") && strings.HasSuffix(candidate, ")")
}

// Get to route root traversing walking up stepping on filler routes
// Preconditions:
// 1. pagePath is a valid page path
func GetRouteRootByWalkingFillerDirs(pagePath string) string {
	routeParts := strings.Split(pagePath, string(os.PathSeparator))
	i := len(routeParts) - 2 // Start from the folder path (not the page.tsx level)
	for ; i > 0; i-- {
		if !isValidFillerDirectory(routeParts[i]) {
			break
		}
	}
	toReturn := strings.Join(routeParts[:i+1], string(os.PathSeparator))
	if toReturn == "" {
		toReturn = string(os.PathSeparator)
	}
	return toReturn
}

// Replace the dynamic route name part of the route name with
// some keyword so that comparision works.
func DynamicRoutePartUnifiedRouteName(name string) string {
	name = GeneralDynamicRouteNameRegex.ReplaceAllString(name, "[slug]")
	name = GeneralDynamicRouteCatchAllNameRegex.ReplaceAllString(name, "[...slug]")
	name = GeneralDynamicRouteOptionalCatchAllNameRegex.ReplaceAllString(name, "[[...slug]]")
	return name
}

func RouteExists(name string, routes []Route, appDir string) (Route, error) {
	var toReturn Route
	// Note here generalized means replacing the "slug" part of [slug] and friends in dynamic
	// route with something universal as [slug] and [foobar] are equivalent for nextJS
	nameGeneralized := DynamicRoutePartUnifiedRouteName(name)
	for _, r := range routes {
		rGeneralized := DynamicRoutePartUnifiedRouteName(RouteFromPagePath(r.PathToPage, appDir))
		// Exact match
		if nameGeneralized == rGeneralized {
			return r, nil
		}
		// TODO
		// Dynamic route but not exact match
		// Meaning, match between CatchAll, OptionalCatchAll, etc.
	}
	return toReturn, fmt.Errorf("route of name %v not found", name)
}

// Get the parent route of the given route
// If no parents exists (for example for route name "/index$"),
// returns the name of the root route, which is ""
//
// Note that this function is different from the GetParent function
// in the frontend where, for example, the parent route name of
// "/products/old/index$" is "/products/old".
//
// Preconditions:
// name is a valid route name
//
// Returns:
// a valid route name
func GetParentRouteName(name string) string {
	// Precondition check
	if err := RouteNameValid(routeName); err != nil {
		panic(err)
	}
	kind, err := RouteTypeFromRouteName(name)
	if err != nil {
		panic(err)
	}
	routeParts := strings.Split(name, "/")
	switch kind {
	case RootRoute:
		// Parent of the root is a root
		return ""
	case StaticRoute:
		// Parent of a static is either filler or root or dynamic
		return strings.Join(routeParts[:len(routeParts)-2], "/")
	case FillerRoute:
		// Parent of a static is either filler or root or dynamic
		return strings.Join(routeParts[:len(routeParts)-1], "/")
	default:
		// Dynamic Route (of any kind)
		return strings.Join(routeParts[:len(routeParts)-2], "/")
	}
}
