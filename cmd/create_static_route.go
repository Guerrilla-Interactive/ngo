package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/Guerrilla-Interactive/ngo/files"
)

// Create static route in the given app directory
// Preconditions:
// name is valid static route name to be created at the location `at`
func createStaticRoute(at string, name string, rawRouteName string) {
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

	CreatedMsg([]string{schemaFilename, queryFilename, pageFilename, previewFilename})
	fmt.Println("")

	// Import the schema into documents.ts
	schemaExportName, err := GetSchemaExportName(name, StaticRoute)
	if err != nil {
		errExit(err)
	}
	err = AddSchemaImportStatement(schemaExportName, schemaFilename)
	if err != nil {
		fmt.Println(err)
	} else {
		schemaDocs, _ := GetSanityDocumentSchemas()
		fmt.Println("Added schema to", schemaDocs)
	}

	// Add the schema into desk structure
	err = AddSchemaToDeskStructure(schemaExportName, StaticRoute)
	if err != nil {
		fmt.Println(err)
	} else {
		deskCustomizationFile, _ := GetSanityDeskCustomozieFileLocation()
		fmt.Println("Added schema to desk structure file", deskCustomizationFile)
	}

	// Append appropriate string to path resolver
	// Note there that 'name' is the last part of the route
	// full name which is the name of the schema.
	// For example for route /chapai/foobar, the schema name is 'foobar'
	err = AddToPathResolver(StaticRoute, name)
	if err != nil {
		fmt.Println(err)
	} else {
		pathResolverFile, _ := GetSanityPathResolverFileLocation()
		fmt.Println("Added schema to path resolver", pathResolverFile)
	}
}
