package main

import "testing"

func TestGotYes(t *testing.T) {
	type test struct {
		query    string
		expected bool
	}
	tests := []test{
		{query: "y", expected: true},
		{query: " y ", expected: true},
		{query: " yes   ", expected: true},
		{query: "YES", expected: true},
		{query: "YES   ", expected: true},
		{query: "  YES   ", expected: true},
		{query: "YeS", expected: true},
		{query: "Ye", expected: false},
		{query: "no", expected: false},
		{query: "Yno", expected: false},
		{query: "Yes n", expected: false},
	}
	for _, test := range tests {
		got := gotYes(test.query)
		if got != test.expected {
			t.Errorf("expected %v got %v for input %v", test.expected, got, test.query)
		}

	}
}
