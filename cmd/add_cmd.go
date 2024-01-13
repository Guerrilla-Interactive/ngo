package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a feature to your nextjs appplication",
	Long: `Add a features to your nextjs application.
For example, you can add a route using the route command. 
Run ngo add route --help for instructions`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("invalid usage. run with --help flag for help")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
