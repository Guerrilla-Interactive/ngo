package cmd

import "testing"

func TestRouteTitleKebabCase(t *testing.T) {
	type TestCase struct {
		name     string
		expected string
	}
	cases := []TestCase{
		{"suman", "suman"},
		{"suman Chapai", "suman-chapai"},
		{"   suman Chapai", "suman-chapai"},
		{"   suman Chapai   ", "suman-chapai"},
		{"suman Chapai   ", "suman-chapai"},
		{" suman   Chapai   ", "suman-chapai"},
		{"  s", "s"},
	}
	for _, testcase := range cases {
		expected := testcase.expected
		got := RouteTitleKebabCase(testcase.name)
		if expected != got {
			t.Errorf("RouteTitleKebabCase(%q) returned %v wanted %v", testcase.name, got, expected)
		}
	}
}

func TestIsRouteNameValid(t *testing.T) {
	type TestCase struct {
		name     string
		expected bool
	}
	cases := []TestCase{
		{"/index/", false}, // must not contain trailing slash
		{"index", false},   // must contain leading slash
		{"/index", true},
		{"/products/index", true},
		{"/products", true},              // Filler route
		{"/products/(something)", false}, // invalid folder name
		{"/(foobar)", false},             // invalid folder name
		{"/[slug]", true},                // invalid folder name
		{"/[index]", false},              // must litrally be "slug" or its friends inside brackets
		{"/products/bogus-filler", true},
		{"/products/[slug]", true},
		{"/products/[...slug]", true},
		{"/products/[[...slug]]", true},
		{"/products/categories/[[...slug]]", true},
		{"/products/categories/[[...index]]", false}, // must literally be "slug" inside brackets
	}
	for _, testcase := range cases {
		expected := testcase.expected
		got := IsRouteNameValid(testcase.name)
		if expected != got {
			t.Errorf("IsRouteNameValid(%q) returned %v wanted %v", testcase.name, got, expected)
		}
	}
}

func TestIsValidFolderName(t *testing.T) {
	type TestCase struct {
		name     string
		expected bool
	}
	cases := []TestCase{
		{"", false}, // expect non-empty
		{"a", true},
		{"ab", true},
		{"abc", true},
		{"(abc)", false},
		{"(abc-def)", false},
		{"(abc-def)abc", false},
		{"abc-def", true},
		{"abc-def-", false},
		{"-abc-def", false},
		{"-abc-def-", false},
	}
	for _, testcase := range cases {
		expected := testcase.expected
		got := IsValidFolderName(testcase.name)
		if expected != got {
			t.Errorf("IsValidFolderName(%q) returned %v wanted %v", testcase.name, got, expected)
		}
	}
}
