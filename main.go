package main

import (
	"fmt"
	"os"
)

func exit(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		exit(err)
	}
	// Check if the directory is empty
	isEmpty, err := IsEmpty(currentDir)
	if err != nil {
		exit(err)
	}
	var wantToContinue string
	if !isEmpty {
		fmt.Printf("%v is not empty\n", currentDir)
		fmt.Printf("Do you still want to continue? ")
		fmt.Scanln(&wantToContinue)
		if !gotYes(wantToContinue) {
			os.Exit(1)
		}
	}
	jsonStr := getSitemapJSON()
	fmt.Println(jsonStr)
}
