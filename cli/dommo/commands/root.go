package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "dommo",
}

func init() {
	rootCmd.AddCommand(
		addDomainCmd,
		deleteDomainCmd,
		getDomainCmd,
		getDomainsFQDNCmd,
		updateDomainCmd,
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
