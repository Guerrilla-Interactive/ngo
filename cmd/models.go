package cmd

type RouteType int

const (
	FillerRoute  RouteType = 0
	StaticRoute  RouteType = 1
	DynamicRoute RouteType = 2
)

type RouteTemplateVariable struct {
	KebabCaseComponentName  string
	PascalCaseComponentName string
	CamelCaseComponentName  string
}
