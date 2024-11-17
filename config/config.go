package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DefaultBackupDir string
	LogDir          string
}

func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &Config{
		DefaultBackupDir: filepath.Join(homeDir, "backups"),
		LogDir:          filepath.Join(homeDir, ".pgbackup", "logs"),
	}, nil
} 