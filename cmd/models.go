package cmd

type RouteType int

const (
	FillerRoute  RouteType = 0
	StaticRoute  RouteType = 1
	DynamicRoute RouteType = 2
)
