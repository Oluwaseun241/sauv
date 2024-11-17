package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/jackc/pgx/v5"
)

type PostgresDB struct {
	conn *pgx.Conn
}

func (pg *PostgresDB) Connect(dbUrl string) error {
	var err error
	pg.conn, err = pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected")
	return nil
}

func (pg *PostgresDB) Backup(dbUrl, destination string) error {
	err := pg.Connect(dbUrl)
	if err != nil {
		return err
	}
	// Check if pg_dump is installed
	if _, err := exec.LookPath("pg_dump"); err != nil {
		return fmt.Errorf("pg_dump is not installed or not in PATH: %v", err)
	}
	cmd := exec.Command("pg_dump", dbUrl, "-Fc", "-f", destination)

	// output(error)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_dump failed: %v\nOutput: %s", err, output)
	}
	log.Printf("Backup successful, saved to %s\n", destination)
	return nil
}
