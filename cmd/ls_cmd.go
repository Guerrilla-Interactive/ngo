package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List routes",
	Long:  `List routes in your project`,
	Run: func(_ *cobra.Command, _ []string) {
		appDir, err := GetAppDirFromWorkingDir()
		if err != nil {
			errExit(err)
		}
		routes := GetRoutes(appDir)
		for _, r := range routes {
			trimmedPath := RouteFromPagePath(r.pathToPage, appDir)
			fmt.Printf("%v - %v\n", r.kind, trimmedPath)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
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
