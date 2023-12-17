package mysqlstore

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestIsConnected(t *testing.T) {
	// Mock a sql.DB for testing purposes
	db := &sql.DB{}
	store := &MySQLTaskStore{db: db}

	if !store.IsConnected() {
		t.Errorf("Expected IsConnected to return true, got false")
	}
}

func TestCreateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection.", err)
	}
	defer db.Close()

	store := &MySQLTaskStore{db: db}
	mock.ExpectExec("INSERT INTO tasks").
		WithArgs("Test Task", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = store.CreateTask("Test Task", []string{"tag1", "tag2"}, time.Now())
	if err != nil {
		t.Errorf("Error was not expected while creating task: %s", err)
	}
}
