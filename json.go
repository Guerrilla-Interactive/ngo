package main

type RouteType int

const (
	FillerRoute  RouteType = 0
	StaticRoute  RouteType = 1
	DynamicRoute RouteType = 2
)

type RootRoute struct {
	ID        string
	Title     string
	Children  []*Route
	SitemapID string
	Type      int
}

type Route struct {
	ID       string
	Title    string
	ParentID string
	Children []*Route
	Type     RouteType
}

type Packages struct {
	// names []string
}

type Sitemap struct {
	ID       string
	Title    string
	Root     *RootRoute
	Packages Packages
}

// func getSitemapStdIn() Sitemap {
// 	// Read from standard input until EOF is found
// 	var sitemap Sitemap
// 	if err := json.NewDecoder(os.Stdin).Decode(&sitemap); err != nil {
// 		log.Fatal(err)
// 	}
// 	return sitemap
// }
