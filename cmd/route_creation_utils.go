package cmd

import (
	"fmt"
	"strings"
)

func AddToPathResolver(routeType RouteType, schemaName string) error {
	magicStringPathResolver := "MAGIC_STRING_SCHEMA_TYPE_TO_PATH_PREFIX\n"
	// Copy the route name from the routeName global variable
	fmt.Println("routename is", routeName)
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
	pathResolverFile, err := GetSanityPathResolverFileLocation()
	if err != nil {
		return err
	}
	AddToFileAfterMagicString(pathResolverFile, magicStringPathResolver, stringToAdd)
	return nil
}

func CreatedMsg(files []string) {
	fmt.Println("Created files:")
	for _, file := range files {
		fmt.Println(file)
	}
}
