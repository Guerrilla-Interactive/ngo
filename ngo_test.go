package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Guerrilla-Interactive/ngo/v2/cmd"
)

func TestRouteTitleToFolderName(t *testing.T) {
	type test struct {
		routeType RouteType
		title     string
		exp       string
	}
	tests := []test{
		{routeType: StaticRoute, title: "About", exp: "about"},
		{routeType: StaticRoute, title: "About Us", exp: "about-us"},
		{routeType: StaticRoute, title: "About   Us", exp: "about-us"},
		{routeType: FillerRoute, title: "Authenticated", exp: "(authenticated)"},
		{routeType: FillerRoute, title: "Authenticated Users", exp: "(authenticated-users)"},
		{routeType: DynamicRoute, title: "Person", exp: "person/[slug]"},
		{routeType: DynamicRoute, title: "Person Name", exp: "person-name/[slug]"},
	}
	for _, test := range tests {
		got := routeTitleToFolderName(test.title, test.routeType)
		if got != test.exp {
			t.Errorf("expected %v got %v for %v\n", test.exp, got, test)
		}
	}
}

func TestCreateFolder(t *testing.T) {
	// Setup
	// Create a tmp directory for playing
	currentDir, err := os.Getwd()
	tmpdir := filepath.Join(currentDir, "_ngo_testdir")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(tmpdir, 0o755)
	if err == os.ErrExist {
		t.Fatalf("folder %v must be empty for testing but already exists", tmpdir)
	} else if err != nil {
		t.Fatal(err)
	}

	type test struct {
		name string
		fail bool
	}
	tests := []test{
		{name: "foo"},
		{name: "bar"},
		{name: "foo/bar"},
		{name: "zoobar/new"}, // Shouldn't fail because createFolder creates intermediate folders as well
	}

	for _, test := range tests {
		_, fail := cmd.CreateFolder(tmpdir, test.name)
		if (fail != nil) != test.fail {
			t.Errorf("expected fail %v but fail %v on creating %v", test.fail, fail != nil, test.name)
		}
	}

	// Ensure that the folders exists
	for _, test := range tests {
		// If failure is expected, then they
		fullPath := filepath.Join(tmpdir, test.name)
		if _, err := os.Stat(fullPath); (err != nil) != test.fail {
			if test.fail {
				t.Errorf("expected %v to not exist but does", fullPath)
			} else {
				t.Errorf("expected %v to exist but doesn't", fullPath)
			}
		}
	}

	// Cleanup
	// Empty contents of tmpdir
	err = os.RemoveAll(tmpdir)
	if err != nil {
		t.Fatalf("err removing tmp dir %v err: %v", tmpdir, err)
	}
}
