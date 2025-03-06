package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sauv/cmd/ui/multiSelect"
	"sauv/cmd/ui/textInput"
	"sauv/internal"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Launch program",
	Run: func(cmd *cobra.Command, args []string) {
		choice := multiSelect.RunSelection()
		if choice == "" {
			fmt.Println("No option selected")
			return
		}

		showAll := choice == "Perform Backup and Migrate"

		values, err := textInput.GetInputs(showAll)
		if err != nil {
			fmt.Printf("Error collecting inputs: %v\n", err)
			return
		}

		oldDBUrl := strings.TrimSpace(values["Old database URL"])
		newDBUrl := strings.TrimSpace(values["New database URL"])
		backupDestination := strings.TrimSpace(values["Backup Destination"])

		switch choice {
		case "Perform Backup":
			if oldDBUrl == "" {
				fmt.Println("You must provide database URL")
				fmt.Println(oldDBUrl)
				return
			}
			if backupDestination == "" {
				fmt.Println("You must provide a backup destination")
				return
			}
			fmt.Println("Performing backup...")
			err := PerformBackup(oldDBUrl, backupDestination)
			if err != nil {
				fmt.Printf("Backup failed: %v\n", err)
				return
			}
			fmt.Println("Backup completed successfully!")

		case "Perform Backup and Migrate":
			if oldDBUrl == "" || newDBUrl == "" {
				fmt.Println("You must provide both old and new database URLs")
				return
			}
			if backupDestination == "" {
				fmt.Println("You must provide a backup destination")
				return
			}
			fmt.Println("Performing backup and migration...")
			err := PerformBackup(oldDBUrl, backupDestination)
			if err != nil {
				fmt.Printf("Backup failed: %v\n", err)
				return
			}
			// TODO: Implement migration
			fmt.Println("Backup completed successfully! Migration in progress")
			err = PerformMigration(oldDBUrl, newDBUrl)
		}
	},
}

func PerformBackup(dbURL, destination string) error {
	if !strings.HasPrefix(dbURL, "postgresql://") {
		return fmt.Errorf("invalid database URL format. Must start with 'postgres://'")
	}

	destDir := filepath.Dir(destination)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Get db type
	dbType := "postgres"

	db, err := internal.GetDatabase(dbType)
	if err != nil {
		fmt.Println("Unsupported database type")
		return err
	}

	err = db.Backup(dbURL, destination)
	if err != nil {
		fmt.Println("Error performing backup:", err)
		return err
	}
	return nil
}

func init() {
	initCmd.Flags().StringVar(&oldDBUrl, "old", "", "Old database URL")
	//backupCmd.Flags().StringVar(&newDBUrl, "new", "", "New database URL")
	rootCmd.AddCommand(initCmd)
}
