package cmd

import "fmt"

// Create dynamic route in the given app directory
func createDynamicRoute(at string, name string, kind DynamicRouteType, rawRouteName string) {
	fmt.Printf("Creating %v at:\n%v\n\n", kind, at)

	// messages := make([]string, 0)
	// mainFolder := filepath.Join(appDir, name)
	// // if the name contains "/" take only the last part as the name of the route
	// routeParts := strings.Split(name, "/")
	//
	// // Note that we are renaming name here to ignore any parent dirs in the name
	// name = routeParts[len(routeParts)-1]
	//
	// slugFolder := filepath.Join(mainFolder, "[slug]")
	// slugCore := filepath.Join(slugFolder, fmt.Sprintf("(%v-slug-core)", name))
	// slugCoreDestination := filepath.Join(slugCore, fmt.Sprintf("(%v-slug-destination)", name))
	// CreatePathAndExitOnFail(slugCoreDestination)
	// // Files: preview and page.tsx and body.tsx
	// slugPreviewFilename := filepath.Join(slugCoreDestination, fmt.Sprintf("%v.slug-preview.tsx", name))
	// CreateFileContents(slugPreviewFilename, files.SlugPreview, name)
	// // slugPageFilename := filepath.Join(slugCoreDestination, "page.tsx")
	// // if catchAllRoute {
	// // 	CreateFileContents(slugPageFilename, files.SlugPageCatchAlll, name)
	// // 	messages = append(messages, fmt.Sprintf("catch all page.tsx: %v", slugPageFilename))
	// // } else {
	// // 	CreateFileContents(slugPageFilename, files.SlugPage, name)
	// // 	messages = append(messages, fmt.Sprintf("page.tsx: %v", slugPageFilename))
	// // }
	// bodyFilename := filepath.Join(slugCoreDestination, fmt.Sprintf("%v.body.tsx", name))
	// CreateFileContents(bodyFilename, files.PageSlugBody, name)
	// messages = append(messages, fmt.Sprintf("page body: %v", bodyFilename))
	//
	// slugCoreServer := filepath.Join(slugCore, fmt.Sprintf("(%v-slug-server)", name))
	// CreatePathAndExitOnFail(slugCoreServer)
	// // Files: slug queries, slug schema
	// slugSchemaFilename := filepath.Join(slugCoreServer, fmt.Sprintf("%v.slug-schema.ts", name))
	// CreateFileContents(slugSchemaFilename, files.SlugSchema, name)
	// messages = append(messages, fmt.Sprintf("slug schema: %v", slugSchemaFilename))
	// slugQueriesFilename := filepath.Join(slugCoreServer, fmt.Sprintf("%v.slug-query.tsx", name))
	// CreateFileContents(slugQueriesFilename, files.SlugQuery, name)
	// messages = append(messages, fmt.Sprintf("slug query: %v", slugQueriesFilename))
	//
	// // Shared utils
	// // Note that here the "shared" refers to shared between the slug and the index page of the route
	// sharedUtilsFolder := filepath.Join(mainFolder, fmt.Sprintf("%v-shared-utils", name))
	// deskStructure := filepath.Join(sharedUtilsFolder, fmt.Sprintf("%v-desk-structure", name))
	// CreatePathAndExitOnFail(deskStructure)
	// // Files: desk-structure.ts
	// deskStructureFilename := filepath.Join(deskStructure, fmt.Sprintf("%v.desk-structure.ts", name))
	// _, err := os.Stat(deskStructureFilename)
	// // Only create the file if it doesn't exist
	// if errors.Is(err, os.ErrNotExist) {
	// 	CreateFileContents(deskStructureFilename, files.SharedDeskStructure, name)
	// 	messages = append(messages, fmt.Sprintf("desk structure: %v", deskStructureFilename))
	// } else if err != nil {
	// 	fmt.Println(err)
	// }
	//
	// sharedQueries := filepath.Join(sharedUtilsFolder, fmt.Sprintf("%v-queries", name))
	// CreatePathAndExitOnFail(sharedQueries)
	// // Files: shared-queries.ts
	// sharedQueriesFilename := filepath.Join(sharedQueries, fmt.Sprintf("%v.shared-queries.ts", name))
	// _, err = os.Stat(sharedQueries)
	// // Only create the file if it doesn't exist
	// if errors.Is(err, os.ErrNotExist) {
	// 	CreateFileContents(sharedQueriesFilename, files.SharedQuery, name)
	// 	messages = append(messages, fmt.Sprintf("desk structure: %v", sharedQueriesFilename))
	// } else if err != nil {
	// 	fmt.Println(err)
	// }
	// printMsg(messages)
}
