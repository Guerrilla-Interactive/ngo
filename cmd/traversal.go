package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func GetAppDirFromWorkingDir() (string, error) {
	var appDir string
	wd, err := os.Getwd()
	if err != nil {
		return appDir, err
	}
	dir, err := getPackageJSONLevelDir(wd)
	if err != nil {
		return appDir, err
	}
	appDir, err = getAppDir(dir)
	if err != nil {
		return appDir, err
	}
	return appDir, nil
}

// Find the direcotry level at which package.json exist. Returns non-nil error
// if cannot find or error exits
func getPackageJSONLevelDir(wd string) (string, error) {
	PACKAGE_JSON := "package.json"
	file := filepath.Join(wd, PACKAGE_JSON)
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		if wd == "/" {
			return wd, fmt.Errorf("cannot find %v in any parent directories", PACKAGE_JSON)
		}
		// Check if package JSON exists in the parent folder
		parentWd := filepath.Dir(wd)
		return getPackageJSONLevelDir(parentWd)
	} else if err != nil {
		// If received a different error, return error
		return wd, err
	}
	return wd, nil
}

// Find the app directory inside the given directory
// Assume dir as the root directory of the project.
// Try dir/app and dir/src/app directory.
func getAppDir(dir string) (string, error) {
	app := filepath.Join(dir, "app")
	_, err := os.Stat(app)
	if err == nil {
		return app, nil
	} else {
		app = filepath.Join(dir, "src/app")
		_, err := os.Stat(app)
		if err != nil {
			errMsg := fmt.Sprintf("No app directory \n%v\n%v", filepath.Join(dir, "app"), filepath.Join(dir, "src/app"))
			return app, errors.New(errMsg)
		}
		return app, nil
	}
}

// Get the sanity documents schema filename string, returns
// error if the file doesn't exist
func GetSanityDocumentSchemas(rootDir string) (string, error) {
	var toReturn string
	toReturn = filepath.Join(rootDir, "sanity/schemas/documents.ts")
	if _, err := os.Stat(toReturn); errors.Is(err, os.ErrNotExist) {
		return toReturn, err
	}
	return toReturn, nil
}

func GetSanityDeskCustomozieFileLocation(rootDir string) (string, error) {
	var toReturn string
	toReturn = filepath.Join(rootDir, "sanity/customize/desk.customize.sanity.tsx")
	if _, err := os.Stat(toReturn); errors.Is(err, os.ErrNotExist) {
		return toReturn, err
	}
	return toReturn, nil
}

func GetSanityPathResolverFileLocation(rootDir string) (string, error) {
	var toReturn string
	toReturn = filepath.Join(rootDir, "sanity/customize/resolve-path.ts")
	if _, err := os.Stat(toReturn); errors.Is(err, os.ErrNotExist) {
		return toReturn, err
	}
	return toReturn, nil
}
