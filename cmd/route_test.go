package cmd

import "testing"

func TestFoldernameToRouteType(t *testing.T) {
	type TestCase struct {
		name string
		kind RouteType
	}
	cases := []TestCase{
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
	type TestCase struct {
		name string
		kind RouteType
		err  error
	}
	var bogusRouteType RouteType
	cases := []TestCase{
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

func TestGetRootRouteByWalkingFillers(t *testing.T) {
	type TestCase struct {
		name     string
		expected string
	}
	cases := []TestCase{
		{"/page.tsx", "/"},
		{"/app/src/(index)/page.tsx", "/app/src"},
		{"/app/src/(index)/(main-destination)/page.tsx", "/app/src"},
		{"/app/src/pieces/(index)/page.tsx", "/app/src/pieces"},
	}
	for _, testcase := range cases {
		got := GetRootRouteByWalkingFillers(testcase.name)
		if testcase.expected != got {
			t.Errorf("GetRootRouteByWalkingFillers(%q) returned %v wanted %v", testcase.name, got, testcase.expected)
		}
	}
}

func TestGetDynamicRouteKindType(t *testing.T) {
	type TestCase struct {
		name     string
		expected DynamicRouteType
	}
	cases := []TestCase{
		{"[slug]", DynamicRoutePrimary},
		{"[...slug]", DynamicRouteCatchAll},
		{"[[...slug]]", DynamicRouteOptionalCatchAll},
	}
	for _, testcase := range cases {
		got, _ := GetDynamicRouteKindType(testcase.name)
		if testcase.expected != got {
			t.Errorf("GetDynamicRouteKindType(%q) returned %v wanted %v", testcase.name, got, testcase.expected)
		}
	}
}
