package cmd

import (
	"errors"
	"fmt"
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
	routeType string // command flag
	addCmd    = &cobra.Command{
		Use:   "add",
		Short: "add a route",
		Long: `Add a route to your next project.
Routes can be static or dynamic. Specify the name of the route using the 
--name flag and type ('dynamic'/'static') using the --type flag).

Creates a static route called about:
ng add --type static --name "/about"

Note that the leading "/" is mandatory and no trailing slash allowed.
Static root route is represented by "/"

To create a dynamic route called categories:
ng add --type dynamic --name "/categories/[slug]"
where the word 'slug' could be replaced with one that you find appropriate

If the dyanmic route is a catchall route, specify it as:
ng add --type dynamic --name "/categories/[...slug]"`,

		Run: func(_ *cobra.Command, _ []string) {
			// Get parsed route type
			r, err := praseRouteType(routeType)
			if err != nil {
				errExit(err)
			}
			if r == FillerRoute {
				errExit(ErrFillerNotImplemented)
			}
			// Get parsed route name (ex. without whitspaces)
			n, err := praseRouteName(routeName)
			if err != nil {
				errExit(err)
			}
			// Check that the route name is valid for the given route type
			// It checks for leading and trailing slashes
			err = AssertRouteNameValid(r, n)
			if err != nil {
				errExit(err)
			}
			createRoute(r, n)
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&routeType, "type", "", "'dyanamic' or 'static'")
	addCmd.Flags().StringVar(&routeName, "name", "", "name of the route")
	addCmd.MarkFlagRequired("type")
	addCmd.MarkFlagRequired("name")
}

// Create a route of given type and name
// If no nextjs project exists (determined by presence of package.json)
// in the current folder or any of the parents then exit with error
// Preconditions:
// 1. r is either a StaticRoute or a DynamicRoute
// 2. name is a valid route name of the appropriate route type (r).
func createRoute(r RouteType, name string) {
	// Precondition check
	if err := AssertRouteNameValid(r, name); err != nil {
		panic(err)
	}
	// Fail if the given route already exists
	appDir, err := GetAppDirFromWorkingDir()
	if err != nil {
		errExit(err)
	}
	routes := GetRoutes(appDir)
	existsRoute := func(candidate string) bool {
		_, err := RouteExists(candidate, routes, appDir)
		return err == nil
	}
	if foundRoute, err := RouteExists(name, routes, appDir); err == nil {
		errExit(fmt.Sprintf("%v already exists",
			RouteFromPagePath(foundRoute.pathToPage, appDir),
		))
	}
	// Remove the leading / is name
	nameLeadinSlashTrimmed := strings.TrimPrefix(name, "/")
	routeParts := strings.Split(nameLeadinSlashTrimmed, "/")
	secondLastRoutePartIndex := len(routeParts) - 2
	// Note that if this is a dynamic route, it contains extra flavor of
	// [...slug] or its variant that needs to be removed
	if r == DynamicRoute {
		secondLastRoutePartIndex--
	}
	var preExistingRouteString string
	for i := 0; i <= secondLastRoutePartIndex; i++ {
		var routeSoFar string
		// Look ahead to check if it's a dynamic route
		// Otherwise the routeSoFar will be incorrectly built
		if i < secondLastRoutePartIndex {
			if IsValidDynamicRouteName(routeParts[i+1]) {
				i++
			}
		}
		routeSoFar = fmt.Sprintf("/%v", strings.Join(routeParts[:i+1], "/"))
		// Check if the route so far exists
		if existsRoute(DynamicRoutePartUnifiedRouteName(routeSoFar)) {
			preExistingRouteString = routeSoFar
		} else {
			break
		}
	}
	// Trim the route name by removing pre-existing route string
	trimmedName := strings.TrimPrefix(name, preExistingRouteString)

	var locationToCreateRoute string
	// If preExistingRoute is "", expected to create at appDir level
	if preExistingRouteString == "" {
		locationToCreateRoute = appDir
	} else {
		// Traverse up from page.tsx of the preExistingRoute page path walking up filler routes
		preExistingRoute, err := RouteExists(preExistingRouteString, routes, appDir)
		locationToCreateRoute = GetRootRouteByWalkingFillers(preExistingRoute.pathToPage)
		if err != nil {
			panic(fmt.Errorf("expected %q to exist, can't find", preExistingRoute.pathToPage))
		}
	}
	routeLocation := filepath.Join(locationToCreateRoute, trimmedName)
	// First we need to find the location at which we can create the route
	// To do that we first find the last parent of the route that doesn't exist
	switch r {
	case DynamicRoute:
		dynamicRouteName := routeParts[len(routeParts)-2]
		dynamicRouteKind, err := GetDynamicRouteKindType(routeParts[len(routeParts)-1])
		if err != nil {
			errExit(err)
		}
		// Create the `routeLocation` folder, along with its parents that don't exist
		// Note that we don't extract this logic for both Static and Dynamic route as
		// fine-graining allows us to create route only after route kind specific error-checking.
		CreatePathAndExitOnFail(routeLocation)
		createDynamicRoute(routeLocation, dynamicRouteName, dynamicRouteKind, routeName)
	case StaticRoute:
		staticRouteName := routeParts[len(routeParts)-1]
		if staticRouteName == "" {
			staticRouteName = "root"
		}
		// Create the `routeLocation` folder, along with its parents that don't exist
		CreatePathAndExitOnFail(routeLocation)
		createStaticRoute(routeLocation, staticRouteName, routeName)
	case FillerRoute:
		errExit(ErrFillerNotImplemented)
	default:
		panic(fmt.Sprintf("unrecognized route - %v\n", r))
	}
}
