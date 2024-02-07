package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/Guerrilla-Interactive/ngo/files"
)

// Create static route in the given app directory
// Preconditions:
// name is valid static route name to be created at the location `at`
func createStaticRoute(at string, name string, rawRouteName string, rootDir string) {
	fmt.Printf("Creating static route at:\n%v\n\n", at)

	schemasAndQueryFolder := filepath.Join(at, fmt.Sprintf("/(index)/(%v-index-core)/(%v-index-server)", name, name))
	pageAndPreviewFolder := filepath.Join(at, fmt.Sprintf("/(index)/(%v-index-core)/(%v-index-destination)", name, name))
	// Create folders
	CreatePathAndExitOnFail(schemasAndQueryFolder)
	CreatePathAndExitOnFail(pageAndPreviewFolder)

	// Schema file
	schemaFilename := filepath.Join(schemasAndQueryFolder, fmt.Sprintf("%v.index-schema.tsx", name))
	CreateFileContents(schemaFilename, files.IndexSchema, name)

	// Query file
	queryFilename := filepath.Join(schemasAndQueryFolder, fmt.Sprintf("%v.index-query.tsx", name))
	CreateFileContents(queryFilename, files.IndexQuery, name)

	// Page file
	pageFilename := filepath.Join(pageAndPreviewFolder, "page.tsx")
	CreateFileContents(pageFilename, files.Page, name)

	// Preview file
	previewFilename := filepath.Join(pageAndPreviewFolder, fmt.Sprintf("%v.index-preview.tsx", name))
	CreateFileContents(previewFilename, files.IndexPreview, name)

	// Page body Component file
	pageBodyComponentFilename := filepath.Join(pageAndPreviewFolder, fmt.Sprintf("%v.index-component.tsx", name))
	CreateFileContents(pageBodyComponentFilename, files.PageIndexBody, name)

	CreatedMsg([]string{
		schemaFilename,
		queryFilename,
		pageFilename,
		previewFilename,
		pageBodyComponentFilename,
	})
	fmt.Println("")

	// Automatically add schema imports, etc. to existing application
	// Note here that name refers to the last part of the route name
	FitNewRouteIntoExistingApp(name, schemaFilename, StaticRoute, rootDir)
}
