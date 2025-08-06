package provider

import (
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func SqlStoreContainer() (*sqlstore.Container, error) {
	sqlDB, err := NewPostgresConnection()
	if err != nil {
		return nil, err
	}
	// defer sqlDB.Close()
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container := sqlstore.NewWithDB(sqlDB, "postgres", dbLog)
	return container, nil
}
