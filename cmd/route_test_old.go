package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFoldernameToRouteType(t *testing.T) {
	type TestCase struct {
		name string
		kind RouteType
	}
	cases := []TestCase{
		{"foobar", StaticRoute},
		{"[slug]", DynamicRoute},
		{"[...index]", DynamicCatchAllRoute},
		{"[[...index]]", DynamicCatchAllOptionalRoute},
	}
	for _, testcase := range cases {
		expected := testcase.kind
		got := FolderNameToRouteType(testcase.name)
		if expected != got {
			t.Errorf("FolderNameToRouteType(%q) returned %v wanted %v", testcase.name, got, expected)
		}
	}
}

func TestRouteTypeByPageTSXPath(t *testing.T) {
	type TestCase struct {
		name string
		kind RouteType
		err  error
	}
	var bogusRouteType RouteType
	cases := []TestCase{
		{filepath.Join("page.tsx"), bogusRouteType, errPagePathInsufficientLength},
		// Note the leading space
		{filepath.Join("foobar", " page.tsx"), bogusRouteType, errPagePathInvalidName},
		// Invalid extension
		{filepath.Join("foobar", "page.ts"), bogusRouteType, errPagePathInvalidName},
		{filepath.Join("foobar", "main.tsx"), bogusRouteType, errPagePathInvalidName},
		{filepath.Join("foobar", "page.tsx"), StaticRoute, nil},
		{filepath.Join("foobar", "page.jsx"), StaticRoute, nil},
		{filepath.Join("(main)", "foobar", "page.tsx"), StaticRoute, nil},
		{filepath.Join("(foobar)", "(xoo)", "xyz", "page.tsx"), StaticRoute, nil},
		{filepath.Join("(foobar)", "[slug]", "xyz", "page.tsx"), StaticRoute, nil},
		{filepath.Join("(foobar)", "[...slug]", "xyz", "page.tsx"), StaticRoute, nil},
		{filepath.Join("(foobar)", "[...slug]", "page.tsx"), DynamicCatchAllRoute, nil},
		{filepath.Join("(foobar)", "xyng", "[...slug]"), bogusRouteType, errPagePathInvalidName},
		// Dynamic route with filler directories
		{filepath.Join("(foobar)", "[slug]", "(xyz)", "page.tsx"), DynamicRoute, nil},
		// Static route at root
		{filepath.Join("(foobar)", "page.jsx"), StaticRoute, nil},
	}
	for _, testcase := range cases {
		expected, expectedErr := testcase.kind, testcase.err
		got, gotErr := RouteTypeByPageTSXPath(testcase.name)
		if expectedErr != gotErr {
			t.Errorf("RouteTypeByPageTSXPath(%q) returned err %v wanted %v", testcase.name, gotErr, expectedErr)
		}
		if expected != got {
			t.Errorf("RouteTypeByPageTSXPath(%q) returned %v wanted %v", testcase.name, got, expected)
		}
	}
}

func TestGetRouteRootByWalkingFillers(t *testing.T) {
	type TestCase struct {
		name     string
		expected string
	}
	cases := []TestCase{
		{filepath.Join("page.tsx"), string(os.PathSeparator)},
		{filepath.Join("app", "src", "(index)", "page.tsx"), filepath.Join("app", "src")},
		{filepath.Join("app", "src", "(index)", "(main-destination)", "page.tsx"), filepath.Join("app", "src")},
		{filepath.Join("app", "src", "pieces", "(index)", "page.tsx"), filepath.Join("app", "src", "pieces")},
	}
	for _, testcase := range cases {
		got := GetRouteRootByWalkingFillerDirs(testcase.name)
		if testcase.expected != got {
			t.Errorf("GetRouteRootByWalkingFillers(%q) returned %v wanted %v", testcase.name, got, testcase.expected)
		}
	}
}

func TestGetParentRouteName(t *testing.T) {
	type TestCase struct {
		name     string
		expected string
	}
	cases := []TestCase{
		{"", ""},
		{"/[slug]", ""},
		{"/[...slug]", ""},
		{"/[[...slug]]", ""},
		{"/(index)", ""},
		{"/products", ""},
		{"/products/categories", "/products"},
		{"/products/categories/(index)", "/products/categories"},
		{"/products/categories/[slug]", "/products/categories"},
		{"/products/categories/[...slug]", "/products/categories"},
		{"/products/categories/[[...slug]]", "/products/categories"},
		{"/products/[slug]/archive", "/products/[slug]"},
	}
	for _, testcase := range cases {
		expected := testcase.expected
		got := GetParentRouteName(testcase.name)
		if expected != got {
			t.Errorf("GetParentRouteName(%q) returned %v wanted %v", testcase.name, got, expected)
		}
	}
}

func TestDynamicRoutePartUnifiedRouteName(t *testing.T) {
	type TestCase struct {
		name     string
		expected string
	}
	cases := []TestCase{
		{"/", "/"},
		{"/foobar", "/foobar"},
		{"/foobar/suman", "/foobar/suman"},
		{"/foobar/[slug]/suman", "/foobar/[slug]/suman"},
		{"/foobar/[...foo]/suman", "/foobar/[...slug]/suman"},
		{"/foobar/[[...foo]]/suman", "/foobar/[[...slug]]/suman"},
	}
	for _, testcase := range cases {
		got := DynamicRoutePartUnifiedRouteName(testcase.name)
		if testcase.expected != got {
			t.Errorf("DynamicRoutePartUnifiedRouteName(%q) returned %v wanted %v", testcase.name, got, testcase.expected)
		}
	}
}
