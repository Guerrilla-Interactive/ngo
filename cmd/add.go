package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Guerrilla-Interactive/ngo/files"
	"github.com/spf13/cobra"
)

// This command is thread unsafe
// because the flags are stored as global variables

// addCmd represents the add command
var (
	routeName             string // command flag
	routeType             string // command flag
	createSchema          bool   // command flag
	createRouteComponents bool   // command flag
	addCmd                = &cobra.Command{
		Use:   "add",
		Short: "add a route",
		Long: `Add a route to your next project.
Routes can be static or dynamic. You should also specify
the name of the route using the -name flag. Example:

ng add -type static -name /about
Creates a static route called about

ng add -type dynamic -name /categories
Creates a dynamic route called categories with pages /categories/[slug]`,
		Run: func(_ *cobra.Command, _ []string) {
			r, err := validateRouteType(routeType)
			if err != nil {
				errExit(err)
			}
			n, err := validateRouteName(routeName)
			if err != nil {
				log.Fatal(err)
			}
			err = createRoute(r, n)
			if err != nil {
				errExit(err)
			}
		},
	}
)

// Create a route of given type and name
// If no nextjs project exists (determined by presence of package.json)
// in the current folder or any of the parents then throw error
func createRoute(r RouteType, name string) error {
	wd, err := os.Getwd()
	if err == nil {
		dir, err := getPackageJSONLevelDir(wd)
		if err != nil {
			return err
		}

		appDir, err := getAppDir(dir)
		if err == nil {
			// TODO
			// Create appropriate files inside the app dir
			switch r {
			case StaticRoute:
				createStaticRoute(appDir, name)
			case DynamicRoute:
				createDynamicRoute(appDir, name)
			}
			return nil
		}
		errExit(err)
	}
	return err
}

// Find the direcotry level at which package.json exist. Returns non-nil error
// if cannot find or error exits
func getPackageJSONLevelDir(wd string) (string, error) {
	PACKAGE_JSON := "package.json"
	file := filepath.Join(wd, PACKAGE_JSON)
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		if wd == "/" {
			return wd, fmt.Errorf("cannot find %v in any parent directories", PACKAGE_JSON)
		}
		// Check if package JSON exists in the parent folder
		parentWd := filepath.Dir(wd)
		return getPackageJSONLevelDir(parentWd)
	} else if err != nil {
		// If received a different error, return error
		return wd, err
	}
	return wd, nil
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&routeType, "type", "t", "", "type dyanamic or type static")
	addCmd.Flags().StringVarP(&routeName, "name", "n", "", "name of the route")
	addCmd.Flags().BoolVarP(&createSchema, "schema", "s", true, "boolean indicating where a schema needs to be created")
	addCmd.MarkFlagRequired("type")
	addCmd.MarkFlagRequired("name")
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

// Find the app directory inside the given directory
// Assume dir as the root directory of the project.
// Try dir/app and dir/src/app directory.
func getAppDir(dir string) (string, error) {
	app := filepath.Join(dir, "app")
	_, err := os.Stat(app)
	if err == nil {
		return app, nil
	} else {
		app = filepath.Join(dir, "src/app")
		_, err := os.Stat(app)
		if err != nil {
			errMsg := fmt.Sprintf("No app directory \n%v\n%v", filepath.Join(dir, "app"), filepath.Join(dir, "src/app"))
			return app, errors.New(errMsg)
		}
		return app, nil
	}
}

// Create static route in the given app directory
func createStaticRoute(appDir string, name string) {
	fmt.Printf("creating static route %s under %s", name, appDir)
	mainFolder := filepath.Join(appDir, name)
	// if the name contains "/" take only the last part as the name of the route
	routeParts := strings.Split(name, "/")
	name = routeParts[len(routeParts)-1]

	// If schema needs to be created
	if createSchema {
		// Note that we create schemas inside (index). We only create the schemas if it doesn't exist already!
		schemasFolder := filepath.Join(mainFolder, fmt.Sprintf("/(index)/(%v-index-core)/(%v-index-server)", name, name))
		CreatePathAndExitOnFail(schemasFolder)
		// Files: schema for the route
		// Add schema to sanity schemas
	}
}

// Create dynamic route in the given app directory
func createDynamicRoute(appDir string, name string) {
	mainFolder := filepath.Join(appDir, name)
	// if the name contains "/" take only the last part as the name of the route
	routeParts := strings.Split(name, "/")

	// Note that we are renaming name here to ignore any parent dirs in the name
	name = routeParts[len(routeParts)-1]

	slugFolder := filepath.Join(mainFolder, "[slug]")
	slugCore := filepath.Join(slugFolder, fmt.Sprintf("%v-slug-core", name))
	slugCoreDestination := filepath.Join(slugCore, fmt.Sprintf("%v-slug-destination", name))
	CreatePathAndExitOnFail(slugCoreDestination)
	// Files: preview and page.tsx
	slugPreviewFilename := filepath.Join(slugCoreDestination, fmt.Sprintf("%v.slug-preview.tsx", name))
	CreateFileContents(slugPreviewFilename, files.SlugPreview, name)
	slugPageFilename := filepath.Join(slugCoreDestination, "page.tsx")
	CreateFileContents(slugPageFilename, files.SlugPage, name)

	slugCoreServer := filepath.Join(slugCore, fmt.Sprintf("%v-slug-server", name))
	CreatePathAndExitOnFail(slugCoreServer)
	// Files: slug queries, slug schema
	if createSchema {
		slugSchemaFilename := filepath.Join(slugCoreServer, fmt.Sprintf("%v.slug-schema.ts", name))
		CreateFileContents(slugSchemaFilename, files.SlugSchema, name)
	}
	slugQueriesFilename := filepath.Join(slugCoreServer, fmt.Sprintf("%v.slug-query.tsx", name))
	CreateFileContents(slugQueriesFilename, files.SlugQuery, name)

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
	} else if err != nil {
		fmt.Println(err)
	}
}
