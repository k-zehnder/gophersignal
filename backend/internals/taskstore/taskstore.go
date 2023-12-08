// taskstore package is the model or data layer
// for our server
package taskstore

import (
	"sync"
	"time"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

// Taskstore is a simple in-memory db of tasks; Taskstore methods are safe
// to call concurrently
type Taskstore struct {
	sync.Mutex

	tasks  map[int]Task
	nextId int
}

func New() *Taskstore {
	return &Taskstore{
		tasks: make(map[int]Task),
	}
}

// CreateTask creates a new task in the store.
func (ts *Taskstore) CreateTask(text string, tags []string, due time.Time) int {
	ts.Lock()
	defer ts.Unlock()

	task := Task{
		Id: ts.nextId,
		Text: text,
		Due: due}
	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)

	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task.Id
}

// // GetTask retrieves a task from the store, by id. If no such id exists, an error is returned.
// func (ts *Taskstore) GetTask(id int) (Task, error)

// // DeleteTask deletes the task with the given id. If no such id exists, an error is returned.
// func (ts *Taskstore) DeleteTask(id int) error

// // DeleteAllTasks deletes all tasks in the store.
// func (ts *Taskstore) DeleteAllTasks() error

// // GetAllTasks returns all the tasks in the store, in arbitrary order.
// func (ts *Taskstore) GetAllTasks() []Task

// // GetTasksByTag returns all the tasks that have the given tag, in arbitrary order
// func (ts *Taskstore) GetTasksByTag(tag string) []Task

// // GetTasksByDueDate
// func (ts *Taskstore) GetTasksByDueDate(year int, month time.Month, day int) []Task
