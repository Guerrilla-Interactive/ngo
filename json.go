package main

import (
	"encoding/json"
	"os"
)

type RootRoute struct {
	ID        string
	Title     string
	Children  []Route
	SitemapID string
	Type      int
}

type Route struct {
	ID       string
	Title    string
	ParentID string
	Children []Route
	Type     int
}

type Packages struct {
	names []string
}

type Sitemap struct {
	ID       string
	Title    string
	Root     RootRoute
	Packages Packages
}

func getSitemapStdIn() Sitemap {
	// Read from standard input until EOF is found
	var sitemap Sitemap
	if err := json.NewDecoder(os.Stdin).Decode(&sitemap); err != nil {
		exit(err)
	}
	return sitemap
}
