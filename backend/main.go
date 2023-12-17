package main

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/k-zehnder/gophersignal/backend/internals/mysqlstore"
)

func main() {
	// exported via ~/.zshrc
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		panic("MYSQL_DSN env variable is not set")
	}
	store, err := mysqlstore.NewMySQLTaskStore(dsn)
	if err != nil {
		panic(err)
	}

	if store.IsConnected() {
		fmt.Println("[x] Successfully connected to the database.")
	} else {
		fmt.Println("[x] Failed to connect to the database.")
	}

	// Create a task
	id, err := store.CreateTask("Test Task", []string{"tag1", "tag2"}, time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
}
