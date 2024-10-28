package db

type MongoDB struct{}

func (mg MongoDB) Connect(conn string) error {
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
