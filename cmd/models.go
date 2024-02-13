package cmd

import (
	"fmt"
	"os"
	"strings"
)

type RouteType int

const (
	FillerRoute  RouteType = 0
	StaticRoute  RouteType = 1
	DynamicRoute RouteType = 2
)

type DynamicRouteType int

const (
	DynamicRoutePrimary          DynamicRouteType = 0
	DynamicRouteCatchAll         DynamicRouteType = 1
	DynamicRouteOptionalCatchAll DynamicRouteType = 2
)

const (
	FillerRouteString  = "Filler"
	StaticRouteString  = "Static"
	DyanmicRouteString = "Dynamic"

	DyanmicRoutePrimaryString          = "dynamic route [slug]"
	DyanmicRouteCatchAllString         = "catchall dynamic route [...slug]"
	DyanmicRouteOptionalCatchAllString = "optional catchall dynamic route [[...slug]]"
)

type RouteTemplateVariable struct {
	KebabCaseComponentName  string
	PascalCaseComponentName string
	CamelCaseComponentName  string
}

type Route struct {
	pathToPage string // Full path (until page.tsx)
	kind       RouteType
}

// Implement the sort interface by RouteLength
type ByRouteLength []Route

func (a ByRouteLength) Len() int      { return len(a) }
func (a ByRouteLength) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByRouteLength) Less(i, j int) bool {
	iLen := len(strings.Split(a[i].pathToPage, "/"))
	jLen := len(strings.Split(a[j].pathToPage, "/"))
	return iLen < jLen
}

func (r RouteType) String() string {
	switch r {
	case FillerRoute:
		return FillerRouteString
	case StaticRoute:
		return StaticRouteString
	default:
		return DyanmicRouteString
	}
}

func (r DynamicRouteType) String() string {
	switch r {
	case DynamicRoutePrimary:
		return DyanmicRoutePrimaryString
	case DynamicRouteCatchAll:
		return DyanmicRouteCatchAllString
	default:
		return DyanmicRouteOptionalCatchAllString
	}
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
	routePartsWithoutFiller := make([]string, 0)
	for _, r := range routeParts {
		if FolderNameToRouteType(r) != FillerRoute {
			routePartsWithoutFiller = append(routePartsWithoutFiller, r)
		}
	}
	result := strings.Join(routePartsWithoutFiller, string(os.PathSeparator))
	// Note trailing slash has to be trimmed before adding a leading slash
	// Remove trailing slash
	result = strings.TrimSuffix(result, string(os.PathSeparator))
	// Add leading slash
	if len(result) == 0 || result[0] != os.PathSeparator {
		result = fmt.Sprintf("%v%v", string(os.PathSeparator), result)
	}
	// Remove any double slashes
	result = strings.ReplaceAll(result, fmt.Sprintf("%v%v", os.PathSeparator, os.PathSeparator), string(os.PathSeparator))
	return result
}
