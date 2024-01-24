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

func TestIsValidDynamicRouteName(t *testing.T) {
	type TestCase struct {
		name     string
		expected bool
	}
	cases := []TestCase{
		{"[foo]", false}, // needs 'slug' literally
		{"bar", false},
		{"/suman", false},
		{"[..slug]", false},
		{"[...slug]", true},
		{"[[...slug]]", true},
		{"[[...slug]", false},
		{"[[..slug]", false},
		{"[[[...slug]]]", false},
		{"[[[[...slug]]]]", false},
		{"[[[...slug chapai]]]", false},
		{"[[...slug chapai]]", false},
	}
	for _, testcase := range cases {
		expected := testcase.expected
		got := IsValidDynamicRouteName(testcase.name)
		if expected != got {
			t.Errorf("IsValidDynamicRouteName(%q) returned %v wanted %v", testcase.name, got, expected)
		}
	}
}
