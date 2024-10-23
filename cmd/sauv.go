package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var oldDBUrl string
var newDBUrl string

var rootCmd = &cobra.Command{
	Use:   "sauv",
	Short: "A database backup and migration utility",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
