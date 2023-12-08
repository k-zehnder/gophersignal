// Basic example of a REST server with several routes, using only the standard library.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/k-zehnder/gophersignal/internals/taskserver"
	"github.com/k-zehnder/gophersignal/internals/taskstore"
)

func main() {
	mux := http.NewServeMux()

	taskstore := taskstore.New()
	taskserver := taskserver.New(taskstore)

	mux.HandleFunc("/task/", taskserver.TaskHandler)

	// Server port set to 3003
	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), mux))


}
