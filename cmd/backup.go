package cmd

import (
	"fmt"
	"sauv/cmd/ui/multiSelect"
	"sauv/cmd/ui/textInput"
	"sauv/internal"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Launch program",
	Run: func(cmd *cobra.Command, args []string) {
		choice := multiSelect.RunSelection()

		values, err := textInput.GetInputs()
		if err != nil {
			fmt.Println("Error collecting inputs:", err)
			return
		}

		oldDBUrl := values["Old database URL"]
		newDBUrl := values["New database URL"]
		backupDestination := values["Backup Destination"]

		switch choice {
		case "Backup":
			if oldDBUrl == "" {
				fmt.Println("You must provide database URL")
				return
			}
			fmt.Println("Performing backup...")
			PerformBackup(oldDBUrl, backupDestination)

		case "Backup and Migrate":
			if oldDBUrl == "" || newDBUrl == "" {
				fmt.Println("You must provide both old and new database URLs")
				return
			}
			fmt.Println("Performing backup and migration...")
		}
	},
}

func PerformBackup(oldDBUrl, backupDestination string) {
	dbType := "postgres"

	db := internal.GetDatabase(dbType)
	if db == nil {
		fmt.Println("Unsupported database type")
		return
	}

	err := db.Backup(oldDBUrl, backupDestination)
	if err != nil {
		fmt.Println("Error performing backup:", err)
		return
	}
}

func init() {
	initCmd.Flags().StringVar(&oldDBUrl, "old", "", "Old database URL")
	//backupCmd.Flags().StringVar(&newDBUrl, "new", "", "New database URL")
	rootCmd.AddCommand(initCmd)
}
