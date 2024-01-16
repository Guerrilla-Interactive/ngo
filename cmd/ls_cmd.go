package cmd

import (
	"fmt"
	"io/fs"
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
				trimmedPath := RouteFromPagePath(r.pathToPage, appDir)
				if list {
					pathFromAppDir := strings.TrimPrefix(r.pathToPage, appDir)
					routeRoot := fmt.Sprintf("app%v", GetRootRouteByWalkingFillers(pathFromAppDir))
					fmt.Printf("%v\t%v\t%v\n", r.kind, trimmedPath, routeRoot)
				} else {
					fmt.Printf("%v\t%v\n", r.kind, trimmedPath)
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
		routeParts := strings.Split(path, "/")
		lastPart := routeParts[len(routeParts)-1]
		return IsValidTerminalPageRouteName(lastPart)
	}
	// Recursively find routes in the appDir
	filepath.WalkDir(appDir, func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if isValidRoutePath(path) {
			routeType, err := RouteTypeByPageTSXPath(path)
			if err != nil {
				return err
			}
			routes = append(routes, Route{path, routeType})
		}
		return nil
	})
	return routes
}
