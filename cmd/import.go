package cmd

import (
	"bytes"
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
