package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/Guerrilla-Interactive/ngo/files"
)

// Create dynamic route in the given app directory
func createDynamicRoute(at string, name string, kind DynamicRouteType, rootDir string) {
	fmt.Printf("Creating %v at:\n%v\n\n", kind, at)

	schemasAndQueryFolder := filepath.Join(at, fmt.Sprintf("/(%v-slug-core)/(%v-slug-server)", name, name))
	pageAndPreviewFolder := filepath.Join(at, fmt.Sprintf("/(%v-slug-core)/(%v-slug-destination)", name, name))
	// Create folders
	CreatePathAndExitOnFail(schemasAndQueryFolder)
	CreatePathAndExitOnFail(pageAndPreviewFolder)

	// Files
	schemaFilename := filepath.Join(schemasAndQueryFolder, fmt.Sprintf("%v.slug-schema.tsx", name))
	queryFilename := filepath.Join(schemasAndQueryFolder, fmt.Sprintf("%v.slug-query.tsx", name))
	pageFilename := filepath.Join(pageAndPreviewFolder, "page.tsx")
	previewFilename := filepath.Join(pageAndPreviewFolder, fmt.Sprintf("%v.slug-preview.tsx", name))
	pageBodyComponentFilename := filepath.Join(pageAndPreviewFolder, fmt.Sprintf("%v.slug-component.tsx", name))

	switch kind {
	case DynamicRoutePrimary:
		CreateFileContents(schemaFilename, files.SlugSchema, name)
		CreateFileContents(queryFilename, files.SlugQuery, name)
		CreateFileContents(pageFilename, files.SlugPage, name) // page.tsx
		CreateFileContents(previewFilename, files.SlugPreview, name)
		CreateFileContents(pageBodyComponentFilename, files.PageSlugBody, name)
	case DynamicRouteCatchAll:
		CreateFileContents(schemaFilename, files.SlugSchemaCatchAll, name)
		CreateFileContents(queryFilename, files.SlugQueryCatchAll, name)
		CreateFileContents(pageFilename, files.SlugPageCatchAll, name) // page.tsx
		CreateFileContents(previewFilename, files.SlugPreviewCatchAll, name)
		CreateFileContents(pageBodyComponentFilename, files.PageSlugBodyCatchAll, name)
	case DynamicRouteOptionalCatchAll:
		CreateFileContents(schemaFilename, files.SlugSchemaCatchAllOptional, name)
		CreateFileContents(queryFilename, files.SlugQueryCatchAllOptional, name)
		CreateFileContents(pageFilename, files.SlugPageCatchAllOptional, name) // page.tsx
		CreateFileContents(previewFilename, files.SlugPreviewCatchAllOptional, name)
		CreateFileContents(pageBodyComponentFilename, files.PageSlugBodyCatchAllOptional, name)
	default:
		errExit(fmt.Sprintf("illegal dynamnic route type got %v", kind))
	}

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
	FitNewRouteIntoExistingApp(name, schemaFilename, DynamicRoute, rootDir)
}
