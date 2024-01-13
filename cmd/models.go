package cmd

import (
	"fmt"
	"strings"
)

type RouteType int

const (
	FillerRoute  RouteType = 0
	StaticRoute  RouteType = 1
	DynamicRoute RouteType = 2
)

const (
	FillerRouteString  = "Filler"
	StaticRouteString  = "Static"
	DyanmicRouteString = "Dynamic"
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

func RouteFromPagePath(path string, appDir string) string {
	trimmedAppDir := strings.TrimPrefix(path, appDir)
	routeParts := strings.Split(trimmedAppDir, "/")
	// Remove last part
	routeParts = routeParts[:len(routeParts)-1]
	routePartsWithoutFiller := make([]string, 0)
	for _, r := range routeParts {
		if FolderNameToRouteType(r) != FillerRoute {
			routePartsWithoutFiller = append(routePartsWithoutFiller, r)
		}
	}
	result := strings.Join(routePartsWithoutFiller, "/")
	// Add leading slash
	if len(result) == 0 || result[0] != '/' {
		result = fmt.Sprintf("/%v", result)
	}
	return result
}
