package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Perform a database backup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Performing backup...")
		PerformBackup(oldDBUrl, newDBUrl)
	},
}

func PerformBackup(oldDBUrl, newDBUrl string) {

}

func init() {
	backupCmd.Flags().StringVar(&oldDBUrl, "old", "", "Old database URL")
	backupCmd.Flags().StringVar(&newDBUrl, "new", "", "New database URL")
	rootCmd.AddCommand(backupCmd)
}
