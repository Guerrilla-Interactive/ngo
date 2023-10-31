package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/Guerrilla-Interactive/ngo/files"
)

type ngo struct {
	rootFolder string
	sitemap    Sitemap
}

// Create folder named `name` under the directory with the path `parent`
// Doesn't create intermediate directories
func createFolder(parent string, name string) (string, error) {
	newName := filepath.Join(parent, name)
	err := os.Mkdir(newName, 0o755)
	return newName, err
}

// Create folder named `name` under the directory with the path `parent` Kills
// the process if any error in creating in the directory like out of space,
// permission error, parent directory doesn't exist etc.
func createFolderAndExitOnFail(parentDir string, name string) string {
	newName, err := createFolder(parentDir, name)
	if err != nil {
		log.Fatal(err)
	}
	return newName
}

// Create file with full path filepath where with data as contents Calls
// `createFile` internally. Returns non-nil error if any error is encoutered
func createFile(filepath string, data []byte) error {
	err := os.WriteFile(filepath, data, 0o644)
	return err
}

// Create file with full path filepath where with data as contents Calls
// `createFile` internally and fails with log.Fatal if any error is encounterd.
func createFileAndExitOnFail(filepath string, data []byte) {
	err := createFile(filepath, data)
	if err != nil {
		log.Fatal(err)
	}
}

// Returns the kebabcase version of the title string
func routeTitleKebabCase(title string) string {
	re := regexp.MustCompile(`\s+`)
	name := strings.ToLower(title)
	// Replace whitespace with -
	name = re.ReplaceAllString(name, "-")
	return name
}

// Return the foldername to use for a route with the provided title and
// RouteType. Transforms spaces into -. Puts square/small/ brackets
// appropriately. Returns the resulting string in kebab-case
func routeTitleToFolderName(title string, routeType RouteType) string {
	name := routeTitleKebabCase(title)
	switch routeType {
	case FillerRoute:
		name = fmt.Sprintf("(%v)", name)
	case DynamicRoute:
		name = "[slug]"
	}
	return name
}

// Recursively create files for a route at given parentDir
func createRouteAt(r *Route, parentDir string) {
	name := routeTitleToFolderName(r.Title, r.Type)
	created := createFolderAndExitOnFail(parentDir, name)
	done := make(chan bool)
	for _, child := range r.Children {
		child := child
		go func() {
			createRouteAt(child, created)
			done <- true
		}()
	}
	// Wait for all route creations to complete
	for i := 0; i < len(r.Children); i++ {
		<-done
	}

	// Create files for each route
	switch r.Type {
	case FillerRoute: // Filler,
		createFillerRouteFilesAt(created, r)
	case StaticRoute: // Static
		createStaticRouteFilesAt(created, r)
	case DynamicRoute: // Dynamic
		createDynamicRouteFilesAt(created, r)
	}
}

// Write contents based on the given template to the given file creating the
// file it it doesn't exist. Note that the template variable is generated using
// generateTemplateVariable function
func createFileContents(filename string, temp *template.Template, r *Route) {
	b := new(bytes.Buffer)
	templateVar := routeTemplateVariable(r.Title)
	if err := temp.Execute(b, templateVar); err != nil {
		log.Fatal(err)
	}
	createFileAndExitOnFail(filename, b.Bytes())
}

// Creates necessary files for a filler route in a given folder
func createFillerRouteFilesAt(folder string, _ *Route) {
	// Create a basic layout.tsx
	// We don't create any files inside the filler
	// file := filepath.Join(folder, "layout.tsx")
	// createFileAndExitOnFail(file, []byte(files.Layout))
}

// Creates necessary files for a static route in a given folder
func createStaticRouteFilesAt(folder string, r *Route) {
	pageNamePrefix := routeTitleKebabCase(r.Title)

	// page.tsx
	file := filepath.Join(folder, "page.tsx")
	createFileContents(file, files.Page, r)

	// page.query.tsx
	query := fmt.Sprintf("%v.query.tsx", pageNamePrefix)
	file = filepath.Join(folder, query)
	createFileContents(file, files.Query, r)

	// page.preview.tsx
	preview := fmt.Sprintf("%v.preview.tsx", pageNamePrefix)
	file = filepath.Join(folder, preview)
	createFileContents(file, files.Preview, r)

	// page.component.tsx
	component := fmt.Sprintf("%v.component.tsx", pageNamePrefix)
	file = filepath.Join(folder, component)
	createFileContents(file, files.Component, r)
}

// Creates necessary files for a dynamic route in a given folder
func createDynamicRouteFilesAt(folder string, r *Route) {
	pageNamePrefix := routeTitleKebabCase(r.Title)

	// page.slug-page.tsx
	pageSlug := fmt.Sprintf("%v.slug-page.tsx", pageNamePrefix)
	file := filepath.Join(folder, pageSlug)
	createFileContents(file, files.SlugPage, r)

	// Create core
	core := createFolderAndExitOnFail(folder, "core")
	serverFolderName := createFolderAndExitOnFail(core, fmt.Sprintf("(%v-server)", pageNamePrefix))
	destinationFolderName := createFolderAndExitOnFail(core, fmt.Sprintf("(%v-destination)", pageNamePrefix))

	// Files inside server
	// page.slug-query.tsx
	file = filepath.Join(serverFolderName, fmt.Sprintf("%v.slug-query.tsx", pageNamePrefix))
	createFileContents(file, files.Query, r)
	// page.slug-schema.ts
	file = filepath.Join(serverFolderName, fmt.Sprintf("%v.slug-schema.ts", pageNamePrefix))
	createFileContents(file, files.QuerySchema, r)

	// Files inside destination
	// page.slug-preview.tsx
	file = filepath.Join(destinationFolderName, fmt.Sprintf("%v.slug-preview.tsx", pageNamePrefix))
	createFileContents(file, files.SlugPreview, r)
	// page.tsx
	file = filepath.Join(destinationFolderName, "page.tsx")
	createFileContents(file, files.Page, r)
}

// Create all necessary files for a given ngo object
func (n *ngo) createFiles() {
	// Create src directory
	src := createFolderAndExitOnFail(n.rootFolder, "src")
	app := createFolderAndExitOnFail(src, "app")
	createFolderAndExitOnFail(src, "components")
	// Create routes inside the app directory
	// TODO root route
	// Create children of root routes
	for _, child := range n.sitemap.Root.Children {
		createRouteAt(child, app)
	}
}
