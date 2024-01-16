package cmd

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

// Create folder named `name` under the directory with the path `parent`
// Create intermediate directories if necessary
func CreateFolder(parent string, name string) (string, error) {
	newName := filepath.Join(parent, name)
	// Create folder, including intermediaries
	err := os.MkdirAll(newName, 0o755)
	return newName, err
}

// Create folder named `name` under the directory with the path `parent` Kills
// the process if any error in creating in the directory like out of space,
// permission error, parent directory doesn't exist etc.
func CreateFolderAndExitOnFail(parentDir string, name string) string {
	newName, err := CreateFolder(parentDir, name)
	if err != nil {
		errExit(err)
	}
	return newName
}

// Create the specified path specified (along with intermediate folders)
// return the path thus created as output. This function exits the program on fail
func CreatePathAndExitOnFail(path string) string {
	err := os.MkdirAll(path, 0o755)
	if err != nil {
		errExit(err)
	}
	return path
}

// Write contents based on the given template to the given file creating the
// file it it doesn't exist. Note that the template variable is generated using
// generateTemplateVariable function
func CreateFileContents(filePath string, temp *template.Template, routeName string) {
	b := new(bytes.Buffer)
	templateVar := GetRouteTemplateVariable(routeName)
	if err := temp.Execute(b, templateVar); err != nil {
		log.Fatal(err)
	}
	err := os.WriteFile(filePath, b.Bytes(), 0o644)
	if err != nil {
		errExit(err)
	}
}
