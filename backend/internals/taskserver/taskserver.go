package taskserver

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"time"

	"github.com/k-zehnder/gophersignal/internals/taskstore"
)

type TaskServer struct {
	store *taskstore.Taskstore
}

// Constructor
func New(store *taskstore.Taskstore) *TaskServer {
	return &TaskServer{store: store}
}

// TaskHandler handles requests to the /task/ endpoint
func (ts *TaskServer) TaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/task/" {
		if req.Method == http.MethodPost {
			ts.createTaskHandler(w, req)
		} else if req.Method == http.MethodGet {
			ts.getAllTasksHandler(w, req)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (ts *TaskServer) createTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling task create at %s\n", req.URL.Path)

	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	type ResponseId struct {
		Id int `json:"id"`
	}

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	var rt RequestTask
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	response := ResponseId{Id: id}
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *TaskServer) getAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all tasks at %s\n", req.URL.Path)
	allTasks := ts.store.GetAllTasks()

	js, err := json.Marshal(allTasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
