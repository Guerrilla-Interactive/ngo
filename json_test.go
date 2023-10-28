package main

import "testing"

func TestRouteType(t *testing.T) {
	type test struct {
		name RouteType
		exp  int
	}
	tests := []test{
		{name: FillerRoute, exp: 0},
		{name: StaticRoute, exp: 1},
		{name: DynamicRoute, exp: 2},
	}
	for _, test := range tests {
		if int(test.name) != test.exp {
			t.Errorf("expected %d got %d", test.exp, test.name)
		}
	}
}
