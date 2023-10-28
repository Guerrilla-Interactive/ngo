package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// Check if the directory is empty
	isEmpty, err := IsEmpty(currentDir)
	if err != nil {
		log.Fatal(err)
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
	sitemap := getSitemapStdIn()
	n := ngo{currentDir, sitemap}
	runPreFilesCreationCommands()
	n.createFiles()
	runPostFilesCreationCommands()
}
