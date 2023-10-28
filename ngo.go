package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
func createFolderAndFail(parentDir string, name string) string {
	newName, err := createFolder(parentDir, name)
	if err != nil {
		log.Fatal(err)
	}
	return newName
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
	created := createFolderAndFail(parentDir, name)
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
	case 0: // Filler,
		createFillerRouteFilesAt(created, r)
	case 1: // Static
		createStaticRouteFilesAt(created, r)
	case 2: // Dynamic
		createDynamicRouteFilesAt(created, r)
	}
}

func createFillerRouteFilesAt(folder string, r *Route) {
}

func createStaticRouteFilesAt(folder string, r *Route) {
}

func createDynamicRouteFilesAt(folder string, r *Route) {
}

func (n *ngo) createFiles() {
	// Create src directory
	src := createFolderAndFail(n.rootFolder, "src")
	app := createFolderAndFail(src, "app")
	createFolderAndFail(src, "components")
	// Create routes inside the app directory
	// TODO root route
	// Create children of root routes
	for _, child := range n.sitemap.Root.Children {
		createRouteAt(child, app)
	}
}
