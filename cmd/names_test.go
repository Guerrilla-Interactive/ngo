package cmd

import "testing"

func TestIsValidDynamicRouteName(t *testing.T) {
	type TestCase struct {
		name     string
		expected bool
	}
	cases := []TestCase{
		{"[foo]", true},
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
