package internal

import (
	"fmt"
	"sauv/internal/db"
)

type Database interface {
	Connect(dbUrl string) error
	Backup(dbUrl, destination string) error
}

func GetDatabase(dbType string) (Database, error) {
	switch dbType {
	case "postgres":
		return &db.PostgresDB{}, nil
	case "mongodb":
		return db.MongoDB{}, nil
	case "sqlite":
		return db.SQLiteDB{}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
