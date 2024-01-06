package cmd

import (
	"fmt"
	"os"
)

// Print the error message and call os.Exit(1)
func errExit(msg interface{}) {
	fmt.Println(msg)
	os.Exit(1)
}
