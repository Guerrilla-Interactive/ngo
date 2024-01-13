package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Guerrilla-Interactive/ngo/files"
	"github.com/spf13/cobra"
)

// Errors
var ErrMultipleFillerRoutesOnAFolder = errors.New("multiple filler routes on a folder")

// This command is thread unsafe
// because the flags are stored as global variables

// routeCmd represents the route subcommand of the add command
var (
	routeName             string // command flag
	routeType             string // command flag
	createSchema          bool   // command flag
	catchAllRoute         bool   // command flag
	createRouteComponents bool   // command flag
	routeCmd              = &cobra.Command{
		Use:   "route",
		Short: "add a route",
		Long: `Add a route to your next project.
Routes can be static or dynamic. Specify the name of the route using the 
--name flag and type ('dynamic'/'static') using the --type flag).

Creates a static route called about:
ng add route --type static --name /about

Note that the leading "/" may be ignored.

To create a dynamic route called categories:
ng add route --type dynamic --name /categories

If the dyanmic route is a catchall route, specify it as:
ng add --type dynamic --name /categories --catchall`,

		Run: func(_ *cobra.Command, _ []string) {
			r, err := validateRouteType(routeType)
			if err != nil {
				errExit(err)
			}
			n, err := validateRouteName(routeName)
			if err != nil {
				errExit(err)
			}
			validateFlagsUsage(r)
			err = createRoute(r, n)
			if err != nil {
				errExit(err)
			}
		},
	}
)

// Validate whether usage of flgas is invalid.
// Note that this function doesn't need to receive flags as arguments
// because it can access it anyway as flags are used as global variables.
func validateFlagsUsage(parsedRouteType RouteType) {
	// Ensure that the catchall flag is not used when the route isn't dynamic
	if parsedRouteType == StaticRoute {
		if catchAllRoute {
			errExit("Catchall flag should only be used with dynamic route")
		}
	}
}

// Create a route of given type and name
// If no nextjs project exists (determined by presence of package.json)
// in the current folder or any of the parents then throw error
func createRoute(r RouteType, name string) error {
	appDir, err := GetAppDirFromWorkingDir()
	if err != nil {
		errExit(err)
	}
	switch r {
	case StaticRoute:
		createStaticRoute(appDir, name)
	case DynamicRoute:
		createDynamicRoute(appDir, name)
	}
	return nil
}

func init() {
	addCmd.AddCommand(routeCmd)
	routeCmd.Flags().StringVar(&routeType, "type", "", "'dyanamic' or 'static'")
	routeCmd.Flags().StringVar(&routeName, "name", "", "name of the route")
	routeCmd.Flags().BoolVar(&createSchema, "schema", true, "boolean indicating where a schema needs to be created")
	routeCmd.Flags().BoolVar(&catchAllRoute, "catchall", false, "boolean indicating whether the dynamic route is catch all")
	routeCmd.MarkFlagRequired("type")
	routeCmd.MarkFlagRequired("name")
}

// Return the route type based on the given candiate string
// also return non nil error when the input is invalid (i.e. no route type
// exists for the given name)
func validateRouteType(candidate string) (RouteType, error) {
	switch candidate {
	case "static":
		return StaticRoute, nil
	case "dynamic":
		return DynamicRoute, nil
	case "filler":
		return FillerRoute, nil
	}
	return StaticRoute, fmt.Errorf("invalid route type. valid types are static, dynamic got %q", candidate)
}

// Return the route name based on the given candiate string
// also return non nil error for when the input is invalid
func validateRouteName(candidate string) (string, error) {
	return RouteTitleKebabCase(candidate), nil
}

// Traverse the current directory walking only on the filler directories
// If any directory contains more than one filler returns error
func LastFillerOnThisFolder(candidate string) (string, error) {
	entries, err := os.ReadDir(candidate)
	if err != nil {
		return "", err
	}
	candidates := make([]string, 0)
	for _, e := range entries {
		if e.IsDir() && FolderNameToRouteType(e.Name()) == FillerRoute {
			candidate := filepath.Join(candidate, e.Name())
			candidates = append(candidates, candidate)
		}
	}
	switch len(candidates) {
	case 0:
		return candidate, nil
	case 1:
		// Recursively find the last filler
		return LastFillerOnThisFolder(candidates[0])
	default:
		return "", ErrMultipleFillerRoutesOnAFolder
	}
}

// Get main folder path given the path to the app directory
// and a routeName string. Note that routeName string may contain
// "/" to separate parent routes. For example:
// To add a page called "values" inside an existing route called "/about"
// the routeName should be specified as "/about/values/", or its variants
// without the first of the or last forward slash.
func GetMainFolderPath(parentDir string, routeName string) (string, error) {
	routeParts := strings.Split(routeName, "/")
	switch len(routeParts) {
	case 0:
		return "", errors.New("no routename provided")
	case 1:
		// Get to the last filler route on this folder
		lastFiller, err := LastFillerOnThisFolder(parentDir)
		if err != nil {
			return "", err
		}
		return lastFiller, nil
	default:
		firstRoute := routeParts[0]
		mainFolderForFirstRoute, err := GetMainFolderPath(parentDir, firstRoute)
		if err != nil {
			return "", err
		}
		// Now we recursively find the main folder for the remaining part of routes
		routeWithoutFirstRoute := strings.Join(routeParts[1:], "/")
		return GetMainFolderPath(mainFolderForFirstRoute, routeWithoutFirstRoute)
	}
}

// Create static route in the given app directory
func createStaticRoute(appDir string, name string) {
	mainFolder, err := GetMainFolderPath(appDir, name)
	if err != nil {
		errExit(err)
	}
	fmt.Println("creating files under", mainFolder)
	// // if the name contains "/" take only the last part as the name of the route
	// routeParts := strings.Split(name, "/")
	// name = routeParts[len(routeParts)-1]
	//
	// // If schema needs to be created
	// if createSchema {
	// 	// Note that we create schemas inside (index). We only create the schemas if it doesn't exist already!
	// 	schemasFolder := filepath.Join(mainFolder, fmt.Sprintf("/(index)/(%v-index-core)/(%v-index-server)", name, name))
	// 	CreatePathAndExitOnFail(schemasFolder)
	// 	// Files: schema for the route
	// 	// Add schema to sanity schemas
	// }
}

// Create dynamic route in the given app directory
func createDynamicRoute(appDir string, name string) {
	messages := make([]string, 0)
	mainFolder := filepath.Join(appDir, name)
	// if the name contains "/" take only the last part as the name of the route
	routeParts := strings.Split(name, "/")

	// Note that we are renaming name here to ignore any parent dirs in the name
	name = routeParts[len(routeParts)-1]

	slugFolder := filepath.Join(mainFolder, "[slug]")
	slugCore := filepath.Join(slugFolder, fmt.Sprintf("(%v-slug-core)", name))
	slugCoreDestination := filepath.Join(slugCore, fmt.Sprintf("(%v-slug-destination)", name))
	CreatePathAndExitOnFail(slugCoreDestination)
	// Files: preview and page.tsx and body.tsx
	slugPreviewFilename := filepath.Join(slugCoreDestination, fmt.Sprintf("%v.slug-preview.tsx", name))
	CreateFileContents(slugPreviewFilename, files.SlugPreview, name)
	slugPageFilename := filepath.Join(slugCoreDestination, "page.tsx")
	if catchAllRoute {
		CreateFileContents(slugPageFilename, files.SlugPageCatchAlll, name)
		messages = append(messages, fmt.Sprintf("catch all page.tsx: %v", slugPageFilename))
	} else {
		CreateFileContents(slugPageFilename, files.SlugPage, name)
		messages = append(messages, fmt.Sprintf("page.tsx: %v", slugPageFilename))
	}
	bodyFilename := filepath.Join(slugCoreDestination, fmt.Sprintf("%v.body.tsx", name))
	CreateFileContents(bodyFilename, files.PageSlugBody, name)
	messages = append(messages, fmt.Sprintf("page body: %v", bodyFilename))

	slugCoreServer := filepath.Join(slugCore, fmt.Sprintf("(%v-slug-server)", name))
	CreatePathAndExitOnFail(slugCoreServer)
	// Files: slug queries, slug schema
	if createSchema {
		slugSchemaFilename := filepath.Join(slugCoreServer, fmt.Sprintf("%v.slug-schema.ts", name))
		CreateFileContents(slugSchemaFilename, files.SlugSchema, name)
		messages = append(messages, fmt.Sprintf("slug schema: %v", slugSchemaFilename))
	}
	slugQueriesFilename := filepath.Join(slugCoreServer, fmt.Sprintf("%v.slug-query.tsx", name))
	CreateFileContents(slugQueriesFilename, files.SlugQuery, name)
	messages = append(messages, fmt.Sprintf("slug query: %v", slugQueriesFilename))

	// Shared utils
	// Note that here the "shared" refers to shared between the slug and the index page of the route
	sharedUtilsFolder := filepath.Join(mainFolder, fmt.Sprintf("%v-shared-utils", name))
	deskStructure := filepath.Join(sharedUtilsFolder, fmt.Sprintf("%v-desk-structure", name))
	CreatePathAndExitOnFail(deskStructure)
	// Files: desk-structure.ts
	deskStructureFilename := filepath.Join(deskStructure, fmt.Sprintf("%v.desk-structure.ts", name))
	_, err := os.Stat(deskStructureFilename)
	// Only create the file if it doesn't exist
	if errors.Is(err, os.ErrNotExist) {
		CreateFileContents(deskStructureFilename, files.SharedDeskStructure, name)
		messages = append(messages, fmt.Sprintf("desk structure: %v", deskStructureFilename))
	} else if err != nil {
		fmt.Println(err)
	}

	sharedQueries := filepath.Join(sharedUtilsFolder, fmt.Sprintf("%v-queries", name))
	CreatePathAndExitOnFail(sharedQueries)
	// Files: shared-queries.ts
	sharedQueriesFilename := filepath.Join(sharedQueries, fmt.Sprintf("%v.shared-queries.ts", name))
	_, err = os.Stat(sharedQueries)
	// Only create the file if it doesn't exist
	if errors.Is(err, os.ErrNotExist) {
		CreateFileContents(sharedQueriesFilename, files.SharedQuery, name)
		messages = append(messages, fmt.Sprintf("desk structure: %v", sharedQueriesFilename))
	} else if err != nil {
		fmt.Println(err)
	}
	printMsg(messages)
}

func printMsg(messages []string) {
	for _, msg := range messages {
		fmt.Println(msg)
	}
}
