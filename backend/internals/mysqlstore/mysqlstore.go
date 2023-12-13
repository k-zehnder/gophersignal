package mysqlstore

import "database/sql"

type MySQLTaskStore struct {
	db *sql.DB
}

func NewMySQLTaskStore(dsn string) (*MySQLTaskStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Create the tasks table if it does not exist
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS tasks (
			id INT UNIQUE AUTO_INCREMENT PRIMARY KEY,
			text VARCHAR(255) NOT NULL,
			tags TEXT,
			due DATETIME
		);
	`
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}

	return &MySQLTaskStore{db: db}, nil
}

func (store *MySQLTaskStore) IsConnected() bool {
	return store.db != nil
}
