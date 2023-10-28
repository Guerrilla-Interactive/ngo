package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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

func createFileAndExitOnFail(filepath string, data string) {
	err := os.WriteFile(filepath, []byte(data), 0o644)
	if err != nil {
		log.Fatal(err)
	}
}

func routeTitleToFolderName(title string, routeType RouteType) string {
	name := strings.ToLower(title)
	// Replace whitespace with -
	name = strings.ReplaceAll(name, " ", "-")
	switch routeType {
	case FillerRoute:
		name = fmt.Sprintf("(%v)", name)
	case DynamicRoute:
		name = "[slug]"
	}
	return name
}

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

func createFillerRouteFilesAt(folder string, _ *Route) {
	// Create a basic layout.tsx
	file := filepath.Join(folder, "layout.tsx")
	createFileAndExitOnFail(file, files.Layout)
}

func createStaticRouteFilesAt(folder string, _ *Route) {
	// page.tsx
	file := filepath.Join(folder, "page.tsx")
	createFileAndExitOnFail(file, "")

	// page.query.tsx
	file = filepath.Join(folder, "page.query.tsx")
	createFileAndExitOnFail(file, "")

	// page.preview.tsx
	file = filepath.Join(folder, "page.preview.tsx")
	createFileAndExitOnFail(file, "")

	// page.component.tsx
	file = filepath.Join(folder, "page.component.tsx")
	createFileAndExitOnFail(file, "")
}

func createDynamicRouteFilesAt(folder string, _ *Route) {
	// page.tsx
	file := filepath.Join(folder, "page.tsx")
	createFileAndExitOnFail(file, "")
}

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
