// cmd/route.go

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	errPagePathInsufficientLength = errors.New("given path to page.tsx has to contain at least one parent directory")
	errPagePathInvalidName        = errors.New("pagename can be page.tsx or page.jsx")
)

// Checks if the given name is valid route terminal name
// Returns true if candidate is 'page.tsx' or 'page.jsx'
func IsValidTerminalPageRouteName(candidate string) bool {
	return candidate == "page.tsx" || candidate == "page.jsx"
}

// FolderNameToRouteType determines the route type based on the folder name.
// Preconditions:
// - name is a valid folder name for Dynamic or a Static route
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

// RouteTypeByPageTSXPath determines the route type based on the path to page.tsx.
// Preconditions:
// - path is a valid full path to page.tsx or page.jsx in a NextJS project with an app directory
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
	// Traverse up from the lowest child ignoring all filler directories
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

// GetRouteRootByWalkingFillerDirs traverses up the directory structure to find the route root,
// skipping filler directories.
// Preconditions:
// - pagePath is a valid page path
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

// DynamicRoutePartUnifiedRouteName replaces the dynamic route name parts with a standard placeholder.
func DynamicRoutePartUnifiedRouteName(name string) string {
	name = GeneralDynamicRouteNameRegex.ReplaceAllString(name, "[slug]")
	name = GeneralDynamicRouteCatchAllNameRegex.ReplaceAllString(name, "[...slug]")
	name = GeneralDynamicRouteOptionalCatchAllNameRegex.ReplaceAllString(name, "[[...slug]]")
	return name
}

// RouteExists checks if a route exists in the given routes list.
func RouteExists(name string, routes []Route, appDir string) (Route, error) {
	var toReturn Route
	// Generalize dynamic route names for comparison
	nameGeneralized := DynamicRoutePartUnifiedRouteName(name)
	for _, r := range routes {
		rGeneralized := DynamicRoutePartUnifiedRouteName(RouteFromPagePath(r.PathToPage, appDir))
		// Exact match
		if nameGeneralized == rGeneralized {
			return r, nil
		}
	}
	return toReturn, fmt.Errorf("route of name %v not found", name)
}

// GetParentRouteName returns the parent route name of the given route.
// If no parent exists, returns the root route "".
// Preconditions:
// - name is a valid route name
func GetParentRouteName(name string) string {
	// Do not replace '/' with OS-specific path separator.

	// Precondition check
	if err := RouteNameValid(name); err != nil {
		panic(err)
	}

	kind, err := RouteTypeFromRouteName(name)
	if err != nil {
		panic(err)
	}

	// Split the route name by '/'
	routeParts := strings.Split(name, "/")

	switch kind {
	case RootRoute:
		// Parent of the root is a root
		return ""
	case StaticRoute, DynamicRoute, DynamicCatchAllRoute, DynamicCatchAllOptionalRoute, FillerRoute:
		// For static and dynamic routes, consider the last non-filler part as the parent.
		// Trim the last part of the route.
		if len(routeParts) > 1 {
			// Exclude the last part to find the parent route.
			parentParts := routeParts[:len(routeParts)-1]

			// Handle routes directly under the root
			if len(parentParts) == 1 && parentParts[0] == "" {
				return "/"
			}

			// Join the parts back together with '/'
			return strings.Join(parentParts, "/")
		}
		return ""
	default:
		panic("unrecognized route types")
	}
}
