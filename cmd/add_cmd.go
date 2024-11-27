package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Errors
var (
	ErrMultipleFillerRoutesOnAFolder = errors.New("multiple filler routes on a folder")
)

// This command is thread unsafe
// because the flags are stored as global variables

// addCmd represents the route subcommand of the add command
var (
	routeName string // command flag
	addCmd    = &cobra.Command{
		Use:   "add",
		Short: "add a route",
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("expected one argument for route name")
			}
			routeName = args[0]
			return nil
		},
		Long: `Add a route to your Next.js project.

Creates a static route called about:
ngo add "/about/(index)"

Creates a static route called values inside about:
ngo add "/about/values/(index)"

Creates a dynamic product route called about:
ngo add "/product/[slug]"

Creates a dynamic catch-all product route called about:
ngo add "/product/[...slug]"

Creates a dynamic optional catch-all product route called about:
ngo add "/product/[[...slug]]"

Creates a filler route:
ngo add "/about/(marketing)"

Note that the leading "/" is mandatory and no trailing slash is allowed.
Static root route is represented by "/(index)"

Dynamic root route is represented by "/[slug]"`,

		Run: func(_ *cobra.Command, _ []string) {
			wd, err := os.Getwd()
			if err != nil {
				errExit(err)
			}
			dir, err := getPackageJSONLevelDir(wd)
			if err != nil {
				errExit(err)
			}
			err = ValidateAndRunAddCommand(routeName, dir)
			if err != nil {
				errExit(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}

// Validates the given route name and route type and runs
// the route add command if the given input is valid assuming
// the projectDir is the root folder of the Next.js project
// where package.json lives.
func ValidateAndRunAddCommand(routeName string, rootDir string) error {
	// Check if the route name is valid
	// Then get the route type from the given name
	_, err := RouteTypeFromRouteName(routeName)
	if err != nil {
		errExit(err)
	}
	// Removed the error exit for filler routes
	return CreateRoute(routeName, rootDir)
}

// Create a route of given type and name
// If no Next.js project exists (determined by presence of package.json)
// in the current folder or any of the parents, then exit with error
// Preconditions:
//  1. name is a valid route name of the appropriate route type (r).
//  2. routeName is either a static route, a dynamic route,
//     or a filler route.
func CreateRoute(routeName, rootDir string) error {
	// Preconditions check
	if err := RouteNameValid(routeName); err != nil {
		panic(err)
	}

	routeKind, err := RouteTypeFromRouteName(routeName)
	if err != nil {
		return err
	}
	if routeKind == RootRoute {
		return nil
	}

	// Get the app directory
	appDir, err := getAppDir(rootDir)
	if err != nil {
		return err
	}

	// Get list of existing routes
	routes := GetRoutes(appDir)
	// Function existsRoute that returns a boolean
	// indicating if the route of a given name already exists
	existsRoute := func(candidate string) bool {
		_, err := RouteExists(candidate, routes, appDir)
		return err == nil
	}

	// Fail if the given route already exists
	if foundRoute, err := RouteExists(routeName, routes, appDir); err == nil {
		errExit(fmt.Sprintf("%v already exists",
			RouteFromPagePath(foundRoute.PathToPage, appDir),
		))
	}

	// We now need:
	// 1. Title for the route (for things like schema name, component name, etc.)
	// 2. Route name, which we already have
	// 3. Route type (static, dynamic, filler)
	// 4. Location to create route

	// Find the closest parent of the route and create
	// the route as its child (by creating the necessary folders to
	// match the route's path)
	parentName := GetParentRouteName(routeName)

	// Check for root route
	for parentName != "" && !existsRoute(parentName) {
		parentName = GetParentRouteName(parentName)
	}

	// Get the location where we need to create the route
	locationToCreateRoute := appDir
	if parentName != "" {
		r, err := RouteExists(parentName, routes, appDir)
		if err != nil {
			return err
		}
		locationToCreateRoute = GetRouteRootByWalkingFillerDirs(r.PathToPage)
	}

	// Create any missing directories in the locationToCreateRoute
	// For example, if creating a static route with name "/about/values/sustainability/(index)",
	// and if "/about" already exists, we need to create the folder "values/sustainability" inside
	// the root of the "about" route.
	foldersToCreate := strings.TrimPrefix(routeName, parentName)
	if routeKind == StaticRoute {
		foldersToCreate = strings.TrimSuffix(foldersToCreate, IndexRouteEnding)
	}
	locationToCreateRoute = filepath.Join(locationToCreateRoute, foldersToCreate)
	// Create all necessary folders
	CreatePathAndExitOnFail(locationToCreateRoute)

	// Create a title for the route that is used for component names and such
	var routeTitle string
	routeParts := strings.Split(routeName, "/")
	if len(routeParts) == 2 {
		routeTitle = "root"
	} else {
		routeTitle = routeParts[len(routeParts)-2]
	}
	if routeKind == StaticRoute {
		createStaticRoute(locationToCreateRoute, routeTitle, routeName, rootDir)
		return nil
	}
	if routeKind == FillerRoute {
		// Call the new function to create a filler route
		err := createFillerRoute(locationToCreateRoute, routeTitle, routeName, rootDir)
		if err != nil {
			return err
		}
		return nil
	}
	// Route kind is Dynamic Route (dynamic, catch-all, etc.)
	createDynamicRoute(locationToCreateRoute, routeTitle, routeKind, rootDir)
	return nil
}

// createFillerRoute creates a filler route with a basic layout.tsx file
func createFillerRoute(path, title, routeName, rootDir string) error {
	// Create the necessary directories
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	// Create a basic layout.tsx file for the filler route
	layoutContent := `import React from 'react';

export default function ` + cases.Title(language.Und).String(title)(title) + `Layout({ children }: { children: React.ReactNode }) {
	return (
		<>
			{children}
		</>
	);
}
`

	// Write the layout.tsx file
	layoutPath := filepath.Join(path, "layout.tsx")
	err = os.WriteFile(layoutPath, []byte(layoutContent), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Filler route %s created successfully at %s\n", routeName, layoutPath)
	return nil
}

// Continue with the rest of your utility functions and handlers...
