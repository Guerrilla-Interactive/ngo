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
		errExit(fmt.Sprintf("%v - route - %v already exists",
			RouteFromPagePath(foundRoute.pathToPage, appDir),
			name,
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

	preExistingRoute, err := RouteExists(preExistingRouteString, routes, appDir)
	if err != nil {
		panic(fmt.Errorf("expected %v to exist, can't find", preExistingRoute.pathToPage))
	}
	var locationToCreateRoute string
	// If preExistingRoute is "", expected to create at appDir level
	if preExistingRouteString == "" {
		locationToCreateRoute = appDir
	} else {
		// Traverse up from page.tsx of the preExistingRoute page path walking up filler routes
		locationToCreateRoute = GetRootRouteByWalkingFillers(preExistingRoute.pathToPage)
	}
	routeLocation := filepath.Join(locationToCreateRoute, trimmedName)
	// First we need to find the location at which we can create the route
	// To do that we first find the last parent of the route that doesn't exist
	switch r {
	case DynamicRoute:
		dynameRouteName := routeParts[len(routeParts)-2]
		dynamicRouteKind, err := GetDynamicRouteKindType(routeParts[len(routeParts)-1])
		if err != nil {
			errExit(err)
		}
		createDynamicRoute(routeLocation, dynameRouteName, dynamicRouteKind)
	case StaticRoute:
		staticRouteName := routeParts[len(routeParts)-1]
		createStaticRoute(routeLocation, staticRouteName)
	case FillerRoute:
		errExit(ErrFillerNotImplemented)
	default:
		panic(fmt.Sprintf("unrecognized route - %v\n", r))
	}
}

// Create static route in the given app directory
// Preconditions:
// name is valid static route name to be created at the location `at`
func createStaticRoute(at string, name string) {
	fmt.Printf("creating static route %v at %v", name, at)
	// if the name contains "/" take only the last part as the name of the route
	routeParts := strings.Split(name, "/")
	name = routeParts[len(routeParts)-1]

	// If schema needs to be created
	// Note that we create schemas inside (index). We only create the schemas if it doesn't exist already!
	schemasFolder := filepath.Join(at, fmt.Sprintf("/(index)/(%v-index-core)/(%v-index-server)", name, name))
	CreatePathAndExitOnFail(schemasFolder)
	// Files: schema for the route
	// Add schema to sanity schemas
}

// Create dynamic route in the given app directory
func createDynamicRoute(at string, name string, kind DynamicRouteType) {
	fmt.Printf("creating %v - %v at %v", kind, name, at)
	// messages := make([]string, 0)
	// mainFolder := filepath.Join(appDir, name)
	// // if the name contains "/" take only the last part as the name of the route
	// routeParts := strings.Split(name, "/")
	//
	// // Note that we are renaming name here to ignore any parent dirs in the name
	// name = routeParts[len(routeParts)-1]
	//
	// slugFolder := filepath.Join(mainFolder, "[slug]")
	// slugCore := filepath.Join(slugFolder, fmt.Sprintf("(%v-slug-core)", name))
	// slugCoreDestination := filepath.Join(slugCore, fmt.Sprintf("(%v-slug-destination)", name))
	// CreatePathAndExitOnFail(slugCoreDestination)
	// // Files: preview and page.tsx and body.tsx
	// slugPreviewFilename := filepath.Join(slugCoreDestination, fmt.Sprintf("%v.slug-preview.tsx", name))
	// CreateFileContents(slugPreviewFilename, files.SlugPreview, name)
	// // slugPageFilename := filepath.Join(slugCoreDestination, "page.tsx")
	// // if catchAllRoute {
	// // 	CreateFileContents(slugPageFilename, files.SlugPageCatchAlll, name)
	// // 	messages = append(messages, fmt.Sprintf("catch all page.tsx: %v", slugPageFilename))
	// // } else {
	// // 	CreateFileContents(slugPageFilename, files.SlugPage, name)
	// // 	messages = append(messages, fmt.Sprintf("page.tsx: %v", slugPageFilename))
	// // }
	// bodyFilename := filepath.Join(slugCoreDestination, fmt.Sprintf("%v.body.tsx", name))
	// CreateFileContents(bodyFilename, files.PageSlugBody, name)
	// messages = append(messages, fmt.Sprintf("page body: %v", bodyFilename))
	//
	// slugCoreServer := filepath.Join(slugCore, fmt.Sprintf("(%v-slug-server)", name))
	// CreatePathAndExitOnFail(slugCoreServer)
	// // Files: slug queries, slug schema
	// slugSchemaFilename := filepath.Join(slugCoreServer, fmt.Sprintf("%v.slug-schema.ts", name))
	// CreateFileContents(slugSchemaFilename, files.SlugSchema, name)
	// messages = append(messages, fmt.Sprintf("slug schema: %v", slugSchemaFilename))
	// slugQueriesFilename := filepath.Join(slugCoreServer, fmt.Sprintf("%v.slug-query.tsx", name))
	// CreateFileContents(slugQueriesFilename, files.SlugQuery, name)
	// messages = append(messages, fmt.Sprintf("slug query: %v", slugQueriesFilename))
	//
	// // Shared utils
	// // Note that here the "shared" refers to shared between the slug and the index page of the route
	// sharedUtilsFolder := filepath.Join(mainFolder, fmt.Sprintf("%v-shared-utils", name))
	// deskStructure := filepath.Join(sharedUtilsFolder, fmt.Sprintf("%v-desk-structure", name))
	// CreatePathAndExitOnFail(deskStructure)
	// // Files: desk-structure.ts
	// deskStructureFilename := filepath.Join(deskStructure, fmt.Sprintf("%v.desk-structure.ts", name))
	// _, err := os.Stat(deskStructureFilename)
	// // Only create the file if it doesn't exist
	// if errors.Is(err, os.ErrNotExist) {
	// 	CreateFileContents(deskStructureFilename, files.SharedDeskStructure, name)
	// 	messages = append(messages, fmt.Sprintf("desk structure: %v", deskStructureFilename))
	// } else if err != nil {
	// 	fmt.Println(err)
	// }
	//
	// sharedQueries := filepath.Join(sharedUtilsFolder, fmt.Sprintf("%v-queries", name))
	// CreatePathAndExitOnFail(sharedQueries)
	// // Files: shared-queries.ts
	// sharedQueriesFilename := filepath.Join(sharedQueries, fmt.Sprintf("%v.shared-queries.ts", name))
	// _, err = os.Stat(sharedQueries)
	// // Only create the file if it doesn't exist
	// if errors.Is(err, os.ErrNotExist) {
	// 	CreateFileContents(sharedQueriesFilename, files.SharedQuery, name)
	// 	messages = append(messages, fmt.Sprintf("desk structure: %v", sharedQueriesFilename))
	// } else if err != nil {
	// 	fmt.Println(err)
	// }
	// printMsg(messages)
}

func printMsg(messages []string) {
	for _, msg := range messages {
		fmt.Println(msg)
	}
}
