package cmd

import (
	"fmt"
	"os"
	"strings"
)

type (
	RouteType        string
	DynamicRouteType RouteType
)

const (
	RootRoute RouteType = "root" // Root route is represented by an empty string
	// Note filler route is the parent of a static route, filler route is not the route
	// with parenthesis. The folder with parenthesis is better called a filler folder/directory
	FillerRoute                  RouteType = "filler"
	StaticRoute                  RouteType = "static"
	DynamicRoute                 RouteType = "dynamic"
	DynamicCatchAllRoute         RouteType = "dynamicCatchAll"
	DynamicCatchAllOptionalRoute RouteType = "dynamicCatchAllOptional"
)

type RouteTemplateVariable struct {
	KebabCaseComponentName  string
	PascalCaseComponentName string
	CamelCaseComponentName  string
}

type Route struct {
	PathToPage string // Full path (until page.tsx), empty string for a filler route
	// This is how the route is represented, for example, in the add command
	RouteRepresentation string
	Kind                RouteType
}

// Implement the sort interface by RouteLength
type ByRouteLength []Route

func (a ByRouteLength) Len() int      { return len(a) }
func (a ByRouteLength) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByRouteLength) Less(i, j int) bool {
	iLen := len(strings.Split(a[i].PathToPage, "/"))
	jLen := len(strings.Split(a[j].PathToPage, "/"))
	return iLen < jLen
}

func (r RouteType) String() string {
	return string(r)
}

// Returns the string representing path for the route
// In case of dynamic routes, representing string such as
// [slug] or [...slug] (for catch-all routes) is present.
// Note that the result contains a leading "/" and no trailing
// slashes.
//
// Preconditions:
//  1. The path is a full path to page.tsx or page.jsx in a nextjs project
func RouteFromPagePath(path string, appDir string) string {
	trimmedAppDir := strings.TrimPrefix(path, appDir)
	routeParts := strings.Split(trimmedAppDir, string(os.PathSeparator))
	// Remove last part
	routeParts = routeParts[:len(routeParts)-1]
	routePartsWithoutFillerDirectories := make([]string, 0)
	for _, r := range routeParts {
		if !isValidFillerDirectory(r) {
			routePartsWithoutFillerDirectories = append(routePartsWithoutFillerDirectories, r)
		}
	}
	result := strings.Join(routePartsWithoutFillerDirectories, string(os.PathSeparator))
	// Note trailing slash has to be trimmed before adding a leading slash
	// Remove trailing slash
	result = strings.TrimSuffix(result, string(os.PathSeparator))
	// Add leading slash
	if len(result) == 0 || result[0] != os.PathSeparator {
		result = fmt.Sprintf("%v%v", string(os.PathSeparator), result)
	}
	// Remove any double slashes
	result = strings.ReplaceAll(result, fmt.Sprintf("%v%v", os.PathSeparator, os.PathSeparator), string(os.PathSeparator))
	// Replace os separator with "/"
	result = strings.Join(strings.Split(result, string(os.PathSeparator)), "/")
	return result
}
