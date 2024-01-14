package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

// Errors
var ErrMultipleFillerRoutesOnAFolder = errors.New("multiple filler routes on a folder")

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
ng add --type static --name /about

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
// in the current folder or any of the parents then throw error
func createRoute(r RouteType, name string) {
	// Fail if the given route already exists
	appDir, err := GetAppDirFromWorkingDir()
	if err != nil {
		errExit(err)
	}
	routes := GetRoutes(appDir)
	switch r {
	case DynamicRoute:
		fmt.Println("pass")
		// for _, r := range routes {
		// 	trimmedPath := RouteFromPagePath(r.pathToPage, appDir)
		// 	routeParts := strings.Split(trimmedPath, "/")
		// 	lastPart := routeParts[len(routeParts)-1]
		// 	switch r.kind {
		// 	case DynamicRoute:
		// 	}
		// }
	case StaticRoute:
		routeStrs := make([]string, 0)
		for _, r := range routes {
			trimmedPath := RouteFromPagePath(r.pathToPage, appDir)
			routeStrs = append(routeStrs, trimmedPath)
		}
		if slices.Contains(routeStrs, name) {
			fmt.Printf("static route %q already exists\n", name)
		} else {
			fmt.Printf("get location to create static route")
			// createStaticRoute(appDir, name)
		}
	case FillerRoute:
		errExit("feature to create a filler route isn't implemented")
	}
}

// // Traverse the current directory walking only on the filler directories
// // If any directory contains more than one filler returns error
//
//	func LastFillerOnThisFolder(candidate string) (string, error) {
//		entries, err := os.ReadDir(candidate)
//		if err != nil {
//			return "", err
//		}
//		candidates := make([]string, 0)
//		for _, e := range entries {
//			if e.IsDir() && FolderNameToRouteType(e.Name()) == FillerRoute {
//				candidate := filepath.Join(candidate, e.Name())
//				candidates = append(candidates, candidate)
//			}
//		}
//		switch len(candidates) {
//		case 0:
//			return candidate, nil
//		case 1:
//			// Recursively find the last filler
//			return LastFillerOnThisFolder(candidates[0])
//		default:
//			return "", ErrMultipleFillerRoutesOnAFolder
//		}
//	}
//
// // Get main folder path given the path to the app directory
// // and a routeName string. Note that routeName string may contain
// // "/" to separate parent routes. For example:
// // To add a page called "values" inside an existing route called "/about"
// // the routeName should be specified as "/about/values/", or its variants
// // without the first of the or last forward slash.
//
//	func GetMainFolderPath(parentDir string, routeName string) (string, error) {
//		routeParts := strings.Split(routeName, "/")
//		switch len(routeParts) {
//		case 0:
//			return "", errors.New("no routename provided")
//		case 1:
//			// Get to the last filler route on this folder
//			lastFiller, err := LastFillerOnThisFolder(parentDir)
//			if err != nil {
//				return "", err
//			}
//			return lastFiller, nil
//		default:
//			firstRoute := routeParts[0]
//			mainFolderForFirstRoute, err := GetMainFolderPath(parentDir, firstRoute)
//			if err != nil {
//				return "", err
//			}
//			// Now we recursively find the main folder for the remaining part of routes
//			routeWithoutFirstRoute := strings.Join(routeParts[1:], "/")
//			return GetMainFolderPath(mainFolderForFirstRoute, routeWithoutFirstRoute)
//		}
//	}

// Create static route in the given app directory
func createStaticRoute(at string, name string) {
	fmt.Println("creating files at", at)
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

// // Create dynamic route in the given app directory
// func createDynamicRoute(appDir string, name string) {
// 	messages := make([]string, 0)
// 	mainFolder := filepath.Join(appDir, name)
// 	// if the name contains "/" take only the last part as the name of the route
// 	routeParts := strings.Split(name, "/")
//
// 	// Note that we are renaming name here to ignore any parent dirs in the name
// 	name = routeParts[len(routeParts)-1]
//
// 	slugFolder := filepath.Join(mainFolder, "[slug]")
// 	slugCore := filepath.Join(slugFolder, fmt.Sprintf("(%v-slug-core)", name))
// 	slugCoreDestination := filepath.Join(slugCore, fmt.Sprintf("(%v-slug-destination)", name))
// 	CreatePathAndExitOnFail(slugCoreDestination)
// 	// Files: preview and page.tsx and body.tsx
// 	slugPreviewFilename := filepath.Join(slugCoreDestination, fmt.Sprintf("%v.slug-preview.tsx", name))
// 	CreateFileContents(slugPreviewFilename, files.SlugPreview, name)
// 	// slugPageFilename := filepath.Join(slugCoreDestination, "page.tsx")
// 	// if catchAllRoute {
// 	// 	CreateFileContents(slugPageFilename, files.SlugPageCatchAlll, name)
// 	// 	messages = append(messages, fmt.Sprintf("catch all page.tsx: %v", slugPageFilename))
// 	// } else {
// 	// 	CreateFileContents(slugPageFilename, files.SlugPage, name)
// 	// 	messages = append(messages, fmt.Sprintf("page.tsx: %v", slugPageFilename))
// 	// }
// 	bodyFilename := filepath.Join(slugCoreDestination, fmt.Sprintf("%v.body.tsx", name))
// 	CreateFileContents(bodyFilename, files.PageSlugBody, name)
// 	messages = append(messages, fmt.Sprintf("page body: %v", bodyFilename))
//
// 	slugCoreServer := filepath.Join(slugCore, fmt.Sprintf("(%v-slug-server)", name))
// 	CreatePathAndExitOnFail(slugCoreServer)
// 	// Files: slug queries, slug schema
// 	slugSchemaFilename := filepath.Join(slugCoreServer, fmt.Sprintf("%v.slug-schema.ts", name))
// 	CreateFileContents(slugSchemaFilename, files.SlugSchema, name)
// 	messages = append(messages, fmt.Sprintf("slug schema: %v", slugSchemaFilename))
// 	slugQueriesFilename := filepath.Join(slugCoreServer, fmt.Sprintf("%v.slug-query.tsx", name))
// 	CreateFileContents(slugQueriesFilename, files.SlugQuery, name)
// 	messages = append(messages, fmt.Sprintf("slug query: %v", slugQueriesFilename))
//
// 	// Shared utils
// 	// Note that here the "shared" refers to shared between the slug and the index page of the route
// 	sharedUtilsFolder := filepath.Join(mainFolder, fmt.Sprintf("%v-shared-utils", name))
// 	deskStructure := filepath.Join(sharedUtilsFolder, fmt.Sprintf("%v-desk-structure", name))
// 	CreatePathAndExitOnFail(deskStructure)
// 	// Files: desk-structure.ts
// 	deskStructureFilename := filepath.Join(deskStructure, fmt.Sprintf("%v.desk-structure.ts", name))
// 	_, err := os.Stat(deskStructureFilename)
// 	// Only create the file if it doesn't exist
// 	if errors.Is(err, os.ErrNotExist) {
// 		CreateFileContents(deskStructureFilename, files.SharedDeskStructure, name)
// 		messages = append(messages, fmt.Sprintf("desk structure: %v", deskStructureFilename))
// 	} else if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	sharedQueries := filepath.Join(sharedUtilsFolder, fmt.Sprintf("%v-queries", name))
// 	CreatePathAndExitOnFail(sharedQueries)
// 	// Files: shared-queries.ts
// 	sharedQueriesFilename := filepath.Join(sharedQueries, fmt.Sprintf("%v.shared-queries.ts", name))
// 	_, err = os.Stat(sharedQueries)
// 	// Only create the file if it doesn't exist
// 	if errors.Is(err, os.ErrNotExist) {
// 		CreateFileContents(sharedQueriesFilename, files.SharedQuery, name)
// 		messages = append(messages, fmt.Sprintf("desk structure: %v", sharedQueriesFilename))
// 	} else if err != nil {
// 		fmt.Println(err)
// 	}
// 	printMsg(messages)
// }
//
// func printMsg(messages []string) {
// 	for _, msg := range messages {
// 		fmt.Println(msg)
// 	}
// }
