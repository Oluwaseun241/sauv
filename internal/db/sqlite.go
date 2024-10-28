package db

type SQLiteDB struct{}

func (sqlite SQLiteDB) Connect(conn string) error {
	return nil
}

func (sqlite SQLiteDB) QueryBackup() error {
	// query the database to get backup
	return nil
}

func (sqlite SQLiteDB) Backup(dbUrl, destination string) error {
	// query the database to get backup
	return nil
}
