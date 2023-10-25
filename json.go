package main

import (
	"fmt"
	"io"
	"strings"
)

func getSitemapJSON() string {
	// Read from standard input until EOF is found
	toReturn := new(strings.Builder)
	for {
		var str string
		_, err := fmt.Scanln(&str)
		if err == io.EOF {
			break
		}
		fmt.Fprintln(toReturn, str)
	}
	return toReturn.String()
}
