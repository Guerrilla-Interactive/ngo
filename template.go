package main

import (
	"fmt"
	"strings"
)

type RouteTemplateVariable struct {
	KebabCaseComponentName  string
	PascalCaseComponentName string
	CamelCaseComponentName  string
}

type DocumentSchemasTemplateVariables struct {
	Schemas []string
}

type PackageJSONTemplateVariables struct {
	PackageName string
}

// Creates a template variable for a given route title. The attributes of this
// variable may be used in the template for a file created for a route.
func routeTemplateVariable(title string) RouteTemplateVariable {
	var v RouteTemplateVariable
	kebab := routeTitleKebabCase(title)

	v.KebabCaseComponentName = kebab
	pascalCase := new(strings.Builder)
	camelCase := new(strings.Builder)
	for i, str := range strings.Split(kebab, "-") {
		if len(str) == 0 {
			continue
		}
		pascal := fmt.Sprintf("%s%s", strings.ToUpper(string(str[0])), str[1:])
		camel := pascal
		if i == 0 {
			camel = str
		}
		pascalCase.WriteString(pascal)
		camelCase.WriteString(camel)
	}
	v.CamelCaseComponentName = camelCase.String()
	v.PascalCaseComponentName = pascalCase.String()
	return v
}
