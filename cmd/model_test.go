package cmd

import "testing"

func TestRoutePathPath(t *testing.T) {
	type TestCase struct {
		path     string
		appDir   string
		epxected string
	}
	cases := []TestCase{
		{
			"/Users/x/Desktop/x/app/(main)/(index)/page.tsx",
			"/Users/x/Desktop/x/app",
			"/",
		},
		{
			"/Users/x/Desktop/x/app/page.tsx",
			"/Users/x/Desktop/x/app",
			"/",
		},
		{
			"/app/page.tsx",
			"/app",
			"/",
		},
		{
			"/app/page.tsx",
			"/app",
			"/",
		},
		{
			"app/page.tsx",
			"app",
			"/",
		},
		{
			"/Users/x/Desktop/x/app/(main-layout)/categories/page.tsx",
			"/Users/x/Desktop/x/app/",
			"/categories",
		},
		{
			"/Users/x/Desktop/x/app/(main-layout)/categories/[slug]/page.tsx",
			"/Users/x/Desktop/x/app",
			"/categories/[slug]",
		},
		{
			"/Users/x/Desktop/x/app/(main-layout)/categories/[...slug]/page.tsx",
			"/Users/x/Desktop/x/app",
			"/categories/[...slug]",
		},
		{
			"/Users/x/Desktop/x/app/(main-layout)/categories/[slug]/archive/page.tsx",
			"/Users/x/Desktop/x/app",
			"/categories/[slug]/archive",
		},
	}
	for _, testcase := range cases {
		expected := testcase.epxected
		got := RouteFromPagePath(testcase.path, testcase.appDir)
		if expected != got {
			t.Errorf("RouteFromPagePath(%q, %q) returned %v wanted %v", testcase.path, testcase.appDir, got, expected)
		}
	}
}
