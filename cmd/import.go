package cmd

import (
	"fmt"
	"os"
	"strings"
)

func AddSchemaImportStatement(path string, name string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	dir, err := getPackageJSONLevelDir(wd)
	if err != nil {
		return err
	}
	importPath := strings.TrimPrefix(path, dir)
	// Import leading slash
	importPath = strings.TrimPrefix(importPath, "/")
	// Remove ts, tsx extension
	importPath = strings.TrimSuffix(importPath, ".tsx")
	importPath = strings.TrimSuffix(importPath, ".ts")
	importString := fmt.Sprintf(`export { %v } from "%v"`, name, importPath)

	documentSchemasFileName, err := GetSanityDocumentSchemas()
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
