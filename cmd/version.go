package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:generate cp ../.version ./
//go:embed .version
var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "See the installed version of ngo",
	Long:  `See the installed version of ngo`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("%v\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
