package taskstore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

type Taskstore struct {
	db *sql.DB
}

func New() *Taskstore {
	// Connect to the database
	dsn := "root:gopher@tcp(localhost:3333)/gopher_api?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	return &Taskstore{db: db}
}

func (ts *Taskstore) CreateTask(text string, tags []string, due time.Time) int {
	// Convert tags slice to JSON for storing in the database
	tagsJSON, _ := json.Marshal(tags)

	// Insert the task into the database
	result, err := ts.db.Exec("INSERT INTO tasks (text, tags, due) VALUES (?, ?, ?)", text, tagsJSON, due)
	if err != nil {
		fmt.Println("Error creating task:", err)
		return 0
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last insert ID:", err)
		return 0
	}

	return int(id)
}

func (ts *Taskstore) GetAllTasks() []Task {
	rows, err := ts.db.Query("SELECT id, text, tags, due FROM tasks")
	if err != nil {
		fmt.Println("Error getting tasks:", err)
		return nil
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var tagsJSON string
		if err := rows.Scan(&task.Id, &task.Text, &tagsJSON, &task.Due); err != nil {
			fmt.Println("Error scanning task:", err)
			continue
		}
		json.Unmarshal([]byte(tagsJSON), &task.Tags)
		tasks = append(tasks, task)
	}

	return tasks
}
