package cmd

import "testing"

func TestFoldernameToRouteType(t *testing.T) {
	type Case struct {
		name string
		kind RouteType
	}
	cases := []Case{
		{"foobar", StaticRoute},
		{"/(main)/foobar", StaticRoute},
		{"/(foobar)/(xoo)/xyz", StaticRoute},
		{"/(foobar)/[slug]/xyz", StaticRoute},
		{"/(foobar)/[...slug]/xyz", StaticRoute},
		{"/(foobar)", FillerRoute},
		{"/(foobar)/[...slug]", DynamicRoute},
		{"/(foobar)/xyng/[...slug]", DynamicRoute},
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
	type Case struct {
		name string
		kind RouteType
		err  error
	}
	var bogusRouteType RouteType
	cases := []Case{
		{"page.tsx", bogusRouteType, errPagePathInsufficientLength},
		// Note the leading space
		{"foobar/ page.tsx", bogusRouteType, errPagePathInvalidName},
		// Invalid extension
		{"foobar/page.ts", bogusRouteType, errPagePathInvalidName},
		{"foobar/main.tsx", bogusRouteType, errPagePathInvalidName},
		{"foobar/page.tsx", StaticRoute, nil},
		{"foobar/page.jsx", StaticRoute, nil},
		{"/(main)/foobar/page.tsx", StaticRoute, nil},
		{"/(foobar)/(xoo)/xyz/page.tsx", StaticRoute, nil},
		{"/(foobar)/[slug]/xyz/page.tsx", StaticRoute, nil},
		{"/(foobar)/[...slug]/xyz/page.tsx", StaticRoute, nil},
		{"/(foobar)/[...slug]/page.tsx", DynamicRoute, nil},
		{"/(foobar)/xyng/[...slug]", bogusRouteType, errPagePathInvalidName},
		// Note that this is a dynamic route with path
		// /[slug]/ and not a filler route,
		// In fact, there's no such thing as filler route,
		// when we are looking as page.tsx. We must thus ignore
		// all filler paths (also called routes) at the lower level
		// to find the actual route type
		{"/(foobar)/[slug]/(xyz)/page.tsx", DynamicRoute, nil},
		// Similarly this is a static route with path "/"
		{"/(foobar)/page.jsx", StaticRoute, nil},
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
