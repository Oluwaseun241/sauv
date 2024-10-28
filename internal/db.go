package internal

import "sauv/internal/db"

type Database interface {
	Connect(dbUrl string) error
	Backup(dbUrl, destination string) error
	QueryBackup() error
}

func GetDatabase(dbType string) Database {
	switch dbType {
	case "postgres":
		return &db.PostgresDB{}
	case "mongodb":
		return db.MongoDB{}
	case "sqlite":
		return db.SQLiteDB{}
	default:
		return nil
	}
}
