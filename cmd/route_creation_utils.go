package cmd

import (
	"fmt"
	"strings"
)

func AddToPathResolver(routeType RouteType, schemaName, rootDir string) error {
	magicStringPathResolver := "MAGIC_STRING_SCHEMA_TYPE_TO_PATH_PREFIX\n"
	// Copy the route name from the routeName global variable
	routeNameCopy := routeName
	if routeType != DynamicRoute && routeType != StaticRoute {
		return fmt.Errorf("route %v not implemented", routeType)
	}
	routeParts := strings.Split(routeName, "/")
	if routeType == DynamicRoute {
		// Drop [...slug] or its friends
		routeNameCopy = strings.Join(routeParts[:len(routeParts)-1], "/")
	}
	stringToAdd := fmt.Sprintf("  %v: '%v',\n", schemaName, routeNameCopy)
	pathResolverFile, err := GetSanityPathResolverFileLocation(rootDir)
	if err != nil {
		return err
	}
	return AddToFileAfterMagicString(pathResolverFile, magicStringPathResolver, stringToAdd)
}

func CreatedMsg(files []string) {
	fmt.Println("Created files:")
	for _, file := range files {
		fmt.Println(file)
	}
}

// Create import statements to import the schema,
// Create appropriate desk structure,
// Update sanity path resolver
func FitNewRouteIntoExistingApp(name, schemaFilename string, kind RouteType, rootDir string) {
	// Import the schema into documents.ts
	schemaExportName, err := GetSchemaExportName(name, kind)
	if err != nil {
		errExit(err)
	}
	err = AddSchemaImportStatement(schemaExportName, schemaFilename, rootDir)
	if err != nil {
		fmt.Println(err)
	} else {
		schemaDocs, _ := GetSanityDocumentSchemas(rootDir)
		fmt.Println("Added schema to", schemaDocs)
	}

	// Add the schema into desk structure
	err = AddSchemaToDeskStructure(schemaExportName, kind, rootDir)
	if err != nil {
		fmt.Println(err)
	} else {
		deskCustomizationFile, _ := GetSanityDeskCustomozieFileLocation(rootDir)
		fmt.Println("Added schema to desk structure file", deskCustomizationFile)
	}

	// Append appropriate string to path resolver
	// Note there that 'name' is the last part of the route
	// full name which is the name of the schema.
	// For example for route /chapai/foobar, the schema name is 'foobar'
	err = AddToPathResolver(kind, name, rootDir)
	if err != nil {
		fmt.Println(err)
	} else {
		pathResolverFile, _ := GetSanityPathResolverFileLocation(rootDir)
		fmt.Println("Added schema to path resolver", pathResolverFile)
	}
}
