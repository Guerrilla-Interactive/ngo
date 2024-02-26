package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var (
	list bool
	// lsCmd represents the ls command
	lsCmd = &cobra.Command{
		Use:   "ls",
		Short: "List routes",
		Long:  `List routes in your project`,
		Run: func(_ *cobra.Command, _ []string) {
			appDir, err := GetAppDirFromWorkingDir()
			if err != nil {
				errExit(err)
			}
			routes := GetRoutes(appDir)
			sort.Sort(ByRouteLength(routes))

			for _, r := range routes {

				// Ignore filler route because it doesn't have path to page.tsx
				if r.Kind == FillerRoute {
					continue
				}

				if list {
					pathFromAppDir := strings.TrimPrefix(r.PathToPage, appDir)
					fmt.Printf("%v\t%v\n", r.RouteRepresentation, pathFromAppDir)
				} else {
					fmt.Printf("%v\n", r.RouteRepresentation)
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().BoolVarP(&list, "list", "l", false, "list the route path")
}

// Returns the list of routes in the given app directory
func GetRoutes(appDir string) []Route {
	routes := make([]Route, 0)
	isValidRoutePath := func(path string) bool {
		routeParts := strings.Split(path, string(os.PathSeparator))
		lastPart := routeParts[len(routeParts)-1]
		return IsValidTerminalPageRouteName(lastPart)
	}
	// Recursively find routes in the appDir
	err := filepath.WalkDir(appDir, func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if isValidRoutePath(path) {
			routeType, err := RouteTypeByPageTSXPath(path)
			if err != nil {
				return err
			}
			routeRootString := RouteFromPagePath(path, appDir)
			if routeType == StaticRoute {
				if routeRootString == "/" {
					routeRootString = "/index"
				} else {
					routeRootString = fmt.Sprintf("%v/index", routeRootString)
				}
			}
			route := Route{
				PathToPage:          path,
				RouteRepresentation: routeRootString,
				Kind:                routeType,
			}
			routes = append(routes, route)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	// From the routes found, add fillers
	// Although we note that users can generate fillers themselves,
	// we are providing it for convinience. Note thought that path to
	// page.tsx isn't present for a filler route.
	fillers := make(map[string]bool)
	for _, r := range routes {
		rRootParts := strings.Split(r.RouteRepresentation, string(os.PathSeparator))
		fillerRouteRoot := strings.Join(rRootParts[:len(rRootParts)-1], string(os.PathSeparator))
		// Root route representation is "", which is also the root filler
		if fillerRouteRoot == "/" {
			fillerRouteRoot = ""
		}
		if _, ok := fillers[fillerRouteRoot]; !ok {
			filler := Route{PathToPage: "", RouteRepresentation: fillerRouteRoot, Kind: FillerRoute}
			routes = append(routes, filler)
			fillers[fillerRouteRoot] = true
		}
	}
	return routes
}
