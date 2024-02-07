package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

// Import sanity schema to document.ts and export from there
func AddSchemaImportStatement(schemaExportName, schemaFilename, rootDir string) error {
	importString, err := SchemaExportImportString(schemaFilename, schemaExportName, rootDir)
	if err != nil {
		return err
	}
	documentSchemasFileName, err := GetSanityDocumentSchemas(rootDir)
	if err != nil {
		return err
	}
	// Write to file, don't create if the file doesn't exist
	f, err := os.OpenFile(documentSchemasFileName, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return (err)
	}
	defer f.Close()
	if _, err = f.WriteString(importString); err != nil {
		return err
	}
	// Append the import string to the document
	return nil
}

func SchemaExportImportString(path, name, rootDir string) (string, error) {
	importPath := strings.TrimPrefix(path, rootDir)
	// Import leading slash
	importPath = strings.TrimPrefix(importPath, "/")
	// Remove ts, tsx extension
	importPath = strings.TrimSuffix(importPath, ".tsx")
	importPath = strings.TrimSuffix(importPath, ".ts")
	return fmt.Sprintf("\nexport { %v } from '%v'", name, importPath), nil
}

func SchemaImportString(name, rootDir string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir, err := getPackageJSONLevelDir(wd)
	if err != nil {
		return "", err
	}
	documentSchemasFileName, err := GetSanityDocumentSchemas(rootDir)
	if err != nil {
		return "", err
	}
	importPath := strings.TrimPrefix(documentSchemasFileName, dir)
	// Import leading slash
	importPath = strings.TrimPrefix(importPath, "/")
	// Remove ts, tsx extension
	importPath = strings.TrimSuffix(importPath, ".tsx")
	importPath = strings.TrimSuffix(importPath, ".ts")
	importPath = "@/" + importPath
	return fmt.Sprintf(`import { %v } from "%v"`, name, importPath), nil
}

// Import sanity schema to document.ts and export from there
func AddSchemaToDeskStructure(schemaExportName string, routeType RouteType, rootDir string) error {
	deskCustomizationFile, err := GetSanityDeskCustomozieFileLocation(rootDir)
	deskStructureImportMagicString := "MAGIC_STRING_CUSTOM_IMPORT\n"
	deskStructureItemMagicString := "MAGIC_STRING_LINE_DESK_STRUCTURES\n"
	if err != nil {
		return err
	}
	var deskStructureItemString string
	switch routeType {
	case StaticRoute:
		deskStructureItemString = fmt.Sprintf("    { type: 'singleton', doc: %v },\n", schemaExportName)
	case DynamicRoute:
		deskStructureItemString = fmt.Sprintf("    { type: 'doc', doc: %v },\n", schemaExportName)
	}
	if deskStructureItemString == "" {
		return fmt.Errorf("invalid route %v for adding schema desk structure", routeType)
	}
	importString, err := SchemaImportString(schemaExportName, rootDir)
	if err != nil {
		return err
	}
	err = AddToFileAfterMagicString(
		deskCustomizationFile,
		deskStructureImportMagicString,
		importString+"\n",
	)
	if err != nil {
		return err
	}
	// Add the schema to desk structure
	err = AddToFileAfterMagicString(
		deskCustomizationFile,
		deskStructureItemMagicString,
		deskStructureItemString,
	)
	if err != nil {
		return err
	}
	return nil
}

// Add stringToAdd to filename after a given magic string. Returns any error encountered.
// The way this function works is that it reads the entire file into memory, splits the file based
// on the given magic string, and joins the file using the magic string
func AddToFileAfterMagicString(filename string, magic string, stringToAdd string) error {
	magicBytes := []byte(magic)
	// Read file into a buffer
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	result := bytes.Split(b, magicBytes)
	if noOfOccurances := len(result) - 1; noOfOccurances != 1 {
		return fmt.Errorf("expected magic string %v to be present once. found %v times", magic, noOfOccurances)
	}
	// Contains magic string exactly once
	// Add the string to right after the magic string
	result[1] = append([]byte(stringToAdd), result[1]...)
	// Expect magic string to end with newline
	if magicBytes[len(magicBytes)-1] != '\n' {
		return fmt.Errorf("expected magic string to end with newline got %q", magic)
	}
	dataToWrite := bytes.Join(result, magicBytes)
	err = os.WriteFile(filename, dataToWrite, 0o644)
	return err
}
