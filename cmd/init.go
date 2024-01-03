package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var (
	templateGitRepo string
	sitemapFileName string
	initCmd         = &cobra.Command{
		Use:   "init",
		Short: "initialize a nextjs project",
		Long: `Use "ngo init" to initalize a nextjs project. 
Choose a custom template via --temp flag.
Optionally provide json sitemap file with the --sitemap flag
to create the required folder structure`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("todo init")
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}

// Returns a slice of git url of template repositories
func getTemplateRepositories() []string {
	// Note that the values in the slice are such that
	// Running git clone xxx should work given that xxx
	// is a value in the slice
	return []string{
		"git@github.com:Guerrilla-Interactive/starter-lun.git",
	}
}
