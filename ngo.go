package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Guerrilla-Interactive/ngo/cmd"
	"github.com/Guerrilla-Interactive/ngo/files"
)

type DocumentSchemasTemplateVariables struct {
	Schemas []string
}

type PackageJSONTemplateVariables struct {
	PackageName string
}

type ngo struct {
	rootFolder string
	sitemap    Sitemap
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

// Return the foldername to use for a route with the provided title and
// RouteType. Transforms spaces into -. Puts square/small/ brackets
// appropriately. Returns the resulting string in kebab-case
func routeTitleToFolderName(title string, routeType RouteType) string {
	name := cmd.RouteTitleKebabCase(title)
	switch routeType {
	case FillerRoute:
		name = fmt.Sprintf("(%v)", name)
	case DynamicRoute:
		name = fmt.Sprintf("%v/[slug]", name)
	}
	return name
}

// Recursively create files for a route at given parentDir
func createRouteAt(r *Route, parentDir string, schemasCh chan<- string) {
	name := routeTitleToFolderName(r.Title, r.Type)
	created := cmd.CreateFolderAndExitOnFail(parentDir, name)
	done := make(chan bool)
	for _, child := range r.Children {
		child := child
		go func() {
			createRouteAt(child, created, schemasCh)
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
		createDynamicRouteFilesAt(created, r, schemasCh)
	}
}

// Write contents based on the given template to the given file creating the
// file it it doesn't exist. Note that the template variable is generated using
// generateTemplateVariable function
func createFileContents(filename string, temp *template.Template, r *Route) {
	b := new(bytes.Buffer)
	templateVar := cmd.GetRouteTemplateVariable(r.Title)
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
	pageNamePrefix := cmd.RouteTitleKebabCase(r.Title)

	// page.tsx
	file := filepath.Join(folder, "page.tsx")
	createFileContents(file, files.Page, r)

	// page.query.tsx
	query := fmt.Sprintf("%v.query.tsx", pageNamePrefix)
	file = filepath.Join(folder, query)
	createFileContents(file, files.IndexQuery, r)

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
func createDynamicRouteFilesAt(folder string, r *Route, schemasCh chan<- string) {
	pageNamePrefix := cmd.RouteTitleKebabCase(r.Title)

	// page.slug-page.tsx
	pageSlug := fmt.Sprintf("%v.slug-page.tsx", pageNamePrefix)
	file := filepath.Join(folder, pageSlug)
	createFileContents(file, files.SlugPage, r)

	// Create core
	core := cmd.CreateFolderAndExitOnFail(folder, "core")
	serverFolderName := cmd.CreateFolderAndExitOnFail(core, fmt.Sprintf("(%v-server)", pageNamePrefix))
	destinationFolderName := cmd.CreateFolderAndExitOnFail(core, fmt.Sprintf("(%v-destination)", pageNamePrefix))

	// Files inside server
	// page.slug-query.tsx
	file = filepath.Join(serverFolderName, fmt.Sprintf("%v.slug-query.tsx", pageNamePrefix))
	createFileContents(file, files.IndexQuery, r)

	// page.slug-schema.ts
	file = filepath.Join(serverFolderName, fmt.Sprintf("%v.slug-schema.ts", pageNamePrefix))
	schemasCh <- file // Send the schema file name to the schemasCh
	createFileContents(file, files.IndexSchema, r)

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
	schemas := make([]string, 0)   // slice of schemas
	schemasCh := make(chan string) // schemasChannel

	// create a goroutine to handle schemas being sent to the channel
	go func() {
		for in := range schemasCh {
			// Replace the first
			in = strings.Replace(in, n.rootFolder, "", 1)
			schemas = append(schemas, in)
		}
	}()

	createPackageJSON(n.rootFolder)
	createFileWithoutTemplateVar(n.rootFolder, ".eslintrc.json", files.Eslintrc)
	createFileWithoutTemplateVar(n.rootFolder, ".gitignore", files.Gitignore)
	createFileWithoutTemplateVar(n.rootFolder, ".npmrc", files.Npmrc)
	createFileWithoutTemplateVar(n.rootFolder, "next.config.mjs", files.NextConfig)
	createFileWithoutTemplateVar(n.rootFolder, "postcss.config.cjs", files.Postcss)
	createFileWithoutTemplateVar(n.rootFolder, "prettier.config.cjs", files.PrettierConfig)
	createFileWithoutTemplateVar(n.rootFolder, "sanity-env.d.ts", files.SanityEnv)
	createFileWithoutTemplateVar(n.rootFolder, "sanity.cli.ts", files.SanityCli)
	createFileWithoutTemplateVar(n.rootFolder, "tailwind.config.cjs", files.Tailwind)
	createFileWithoutTemplateVar(n.rootFolder, "tsconfig.json", files.TSConfigJSON)

	// Create src directory
	src := cmd.CreateFolderAndExitOnFail(n.rootFolder, "src")
	app := cmd.CreateFolderAndExitOnFail(src, "app")

	// Sanity related directories
	sanity := cmd.CreateFolderAndExitOnFail(n.rootFolder, "sanity")
	sanitySchemas := cmd.CreateFolderAndExitOnFail(sanity, "schemas")
	cmd.CreateFolderAndExitOnFail(src, "components")

	// Create routes inside the app directory
	// TODO root route
	// Create children of root routes
	for _, child := range n.sitemap.Root.Children {
		createRouteAt(child, app, schemasCh)
	}

	// Sanity schemas file
	// Create document index.ts for schemas
	fileDocIndexSchemas := filepath.Join(sanitySchemas, "index.ts")
	b := new(bytes.Buffer)
	if err := files.DocumentsSchemaQuery.Execute(b, DocumentSchemasTemplateVariables{schemas}); err != nil {
		log.Fatal(err)
	}
	createFileAndExitOnFail(fileDocIndexSchemas, b.Bytes())
}

func createPackageJSON(location string) {
	names := strings.Split(location, "/")
	name := names[len(names)-1]
	templateVar := PackageJSONTemplateVariables{name}
	filePath := filepath.Join(location, "package.json")
	b := new(bytes.Buffer)
	if err := files.PackageJSON.Execute(b, templateVar); err != nil {
		log.Fatal(err)
	}
	createFileAndExitOnFail(filePath, b.Bytes())
}

func createFileWithoutTemplateVar(location string, filename string, temp *template.Template) {
	filePath := filepath.Join(location, filename)
	b := new(bytes.Buffer)
	if err := temp.Execute(b, nil); err != nil {
		log.Fatal(err)
	}
	createFileAndExitOnFail(filePath, b.Bytes())
}
