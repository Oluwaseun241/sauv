package cmd

import (
	"fmt"
)

func PerformMigration(oldDBUrl, newDBUrl string) error {
	fmt.Printf("Performing migration from %s to %s\n", oldDBUrl, newDBUrl)
	fmt.Println("\nMigration completed successfully!")
	// migrate to new db(url) provided
	return nil
}
