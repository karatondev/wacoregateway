package provider

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/faisolarifin/wacoregateway/util"
	_ "github.com/lib/pq"
)

func NewPostgresConnection() (*sql.DB, error) {
	cfg := util.Configuration.Postgres

	// Create the connection string
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		strings.Join(cfg.Options, "&"),
	)

	// Open the connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	// an db without sqlx.DB instance
	return db, nil
}
