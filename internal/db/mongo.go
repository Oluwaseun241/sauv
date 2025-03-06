package db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct{}

func (mg MongoDB) Connect(dbUrl string) error {
	client, _ := mongo.Connect(options.Client().ApplyURI(dbUrl))
	fmt.Println(client)
	return nil
}

func (mg MongoDB) QueryBackup() error {
	// query the database to get backup
	return nil
}

func (mg MongoDB) Backup(dbUrl, destination string) error {
	// query the database to get backup
	return nil
}
