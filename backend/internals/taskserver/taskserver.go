package taskserver

import (
	"encoding/json"
	"fmt"
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

// Declaring method TaskHandler on TaskServer struct
func (ts *TaskServer) TaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/task/" {
		// Request is plain "/task/", without trailing ID.
		if req.Method == http.MethodPost {
			ts.createTaskHandler(w, req)
		}
	}
}

func (ts *TaskServer) createTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling task create at %s\n", req.URL.Path)

	// Types used internall in this handler to (de-)serialize the request and
	// response to/from JSON
	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	type ResponseId struct {
		Id int `json:"id"`
	}

	// Enforce a JSON Content-Type
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rt RequestTask
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	fmt.Println("[x] id:", id)

}
