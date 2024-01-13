package cmd

import (
	"errors"
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

func FolderNameToRouteType(name string) RouteType {
	// If foldername is the path to the folder, extract
	// just the last part so we get the folder name instead
	routeParts := strings.Split(name, "/")
	lastPart := routeParts[len(routeParts)-1]
	if strings.HasPrefix(lastPart, "(") && strings.HasSuffix(lastPart, ")") {
		return FillerRoute
	} else if strings.HasPrefix(lastPart, "[") && strings.HasSuffix(lastPart, "]") {
		return DynamicRoute
	} else {
		return StaticRoute
	}
}

func RouteTypeByPageTSXPath(path string) (RouteType, error) {
	var kind RouteType
	routeParts := strings.Split(path, "/")
	lastPart := routeParts[len(routeParts)-1]
	if len(routeParts) < 2 {
		return kind, errPagePathInsufficientLength
	}
	if !IsValidTerminalPageRouteName(lastPart) {
		return kind, errPagePathInvalidName
	}
	// Now we traverse up from the lowest child ignoring all filler routes
	for i := len(routeParts) - 2; i >= 0; i-- {
		name := routeParts[i]
		if kind := FolderNameToRouteType(name); kind != FillerRoute {
			return kind, nil
		}
	}
	return StaticRoute, nil
}
