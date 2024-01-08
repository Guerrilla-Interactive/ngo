package cmd

import "testing"

func TestRouteVariableTemplate(t *testing.T) {
	type Test struct {
		title string
		exp   RouteTemplateVariable
	}
	tests := []Test{
		{"Test", RouteTemplateVariable{"test", "Test", "test"}},
		{"Test Name", RouteTemplateVariable{"test-name", "TestName", "testName"}},
		{"Test Name   Multiple space", RouteTemplateVariable{"test-name-multiple-space", "TestNameMultipleSpace", "testNameMultipleSpace"}},
	}
	for _, test := range tests {
		got := GetRouteTemplateVariable(test.title)
		if got != test.exp {
			t.Errorf("expected variable %v got %v", test.exp, got)
		}
	}
}
