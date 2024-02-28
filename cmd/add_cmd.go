package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Errors
var (
	ErrMultipleFillerRoutesOnAFolder = errors.New("multiple filler routes on a folder")
	ErrFillerNotImplemented          = errors.New("ability to add filler route is not implemented")
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
		Long: `Add a route to your next project.

Creates a static route called about:
ng add "/about/index$"

Creates a static route called values inside about:
ng add "/about/values/index$"

Creates a dynamic product route called about:
ng add "/product/[slug]"

Creates a dynamic catchall product route called about:
ng add "/product/[...slug]"

Creates a dynamic catchall optional product route called about:
ng add "/product/[[...slug]]"

Note that the leading "/" is mandatory and no trailing slash allowed.
Static root route is represented by "/index$"

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
// the the route add command if the given input is valid assuming
// the projectDir is the root folder of the NextJS project.
// where package.JSON lives.
func ValidateAndRunAddCommand(routeName string, rootDir string) error {
	// Check if the route name is valid
	// Then get the route type from the given name
	r, err := RouteTypeFromRouteName(routeName)
	if err != nil {
		errExit(err)
	}
	if r == FillerRoute {
		errExit(ErrFillerNotImplemented)
	}
	return CreateRoute(routeName, rootDir)
}

// Create a route of given type and name
// If no nextjs project exists (determined by presence of package.json)
// in the current folder or any of the parents then exit with error
// Preconditions:
// 1. name is a valid route name of the appropriate route type (r).
// 2. routeName is a either static route or one of the dyanamic routes (dynamic,
// dynamic catch all or dynamic catch all optional)
func CreateRoute(routeName, rootDir string) error {
	// Preconditions check
	if err := RouteNameValid(routeName); err != nil {
		panic(err)
	}

	routeKind, err := RouteTypeFromRouteName(routeName)
	if err != nil {
		return err
	}
	if routeKind == FillerRoute {
		panic("creating filler route is unsupported")
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
	// Anonymous function existsRoute that returns a boolean
	// indicating if the route of a given name already exists
	existsRoute := func(candidate string) bool {
		_, err := RouteExists(candidate, routes, appDir)
		return err == nil
	}

	// Fail if the given route already exists

	// Note that we are not using the existsRoute function here because
	// we also need the result of where the route exists
	// (i.e the first thing that RouteExists returns)
	if foundRoute, err := RouteExists(routeName, routes, appDir); err == nil {
		errExit(fmt.Sprintf("%v already exists",
			RouteFromPagePath(foundRoute.PathToPage, appDir),
		))
	}

	// We now need
	// 1. Title for the route (for things like schema name, component name, etc...)
	// 2. Route name, which we already have
	// 3. Route type (static or dynamic)
	// 4. Location to create route

	// We find the closest parent of the route and create
	// the route as its child (by creating the necessary folders to
	// match the route's path)
	parentName := GetParentRouteName(routeName)

	// Note that in several places in the code below, we are conditionally
	// checking for parent == "" which tells whether the parent is a root route
	// For example, when creating /index$ static route, its parent is ""

	// Find the first parent that exists or the root route if no parent exists
	for parentName != "" && !existsRoute(parentName) {
		parentName = GetParentRouteName(parentName)
	}

	// We now get the location where we need to create the route
	locationToCreateRoute := appDir
	if parentName != "" {
		r, err := RouteExists(parentName, routes, appDir)
		if err != nil {
			return err
		}
		locationToCreateRoute = GetRouteRootByWalkingFillerDirs(r.PathToPage)
	}

	// Create any missing directories in the locationToCreateRoute
	// For example, if creating a static route with name "/about/values/sustainability/index$",
	// if the "/about" already exists, we need to create th folder "values/sustainability" inside
	// the root of the "about" route (by route we mean, first parent of about's page.tsx that's
	// not a filler directory)
	foldersToCreate := strings.TrimPrefix(routeName, parentName)
	if routeKind == StaticRoute {
		foldersToCreate = strings.TrimSuffix(foldersToCreate, IndexRouteEnding)
	}
	// This is the filepath level where we can create page.tsx files
	// Note though that we might, in reality, create the page.tsx file inside
	// a filler directory in this path.
	locationToCreateRoute = filepath.Join(locationToCreateRoute, foldersToCreate)
	// We create all necessary foldrs by calling this function
	CreatePathAndExitOnFail(locationToCreateRoute)

	// Now we create a title for the route that is used for component names and such
	var routeTitle string
	routeParts := strings.Split(routeName, string(os.PathSeparator))
	if len(routeParts) == 2 {
		routeTitle = "root"
	} else {
		routeTitle = routeParts[len(routeParts)-2]
	}
	if routeKind == StaticRoute {
		createStaticRoute(locationToCreateRoute, routeTitle, routeName, rootDir)
		return nil
	}
	// Route kind is Dyanmic Route (of any variant: catchall, catchall optional, etc...)
	createDynamicRoute(locationToCreateRoute, routeTitle, routeKind, rootDir)
	return nil
}
