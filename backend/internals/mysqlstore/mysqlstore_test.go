package mysqlstore

import (
	"database/sql"
	"testing"
)

func TestIsConnected(t *testing.T) {
	// Mack a sql.DB for testing purposes
	db := &sql.DB{}
	store := &MySQLTaskStore{db: db}

	if !store.IsConnected() {
		t.Errorf("Expected IsConnected to return true, got false")
	}
}
