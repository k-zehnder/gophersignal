package main

import (
	"log"
	"net/http"

	"github.com/k-zehnder/gophersignal/config"
	"github.com/k-zehnder/gophersignal/internal/api/controller"
	_ "github.com/k-zehnder/gophersignal/internal/api/docs"
	"github.com/k-zehnder/gophersignal/internal/api/router"
	"github.com/k-zehnder/gophersignal/internal/store"
)

func main() {
	sqlconnection := config.NewConnection()
	store := store.NewStore(sqlconnection)
	controller := controller.NewController(store)
	routes := router.NewRouter(controller)

	server := &http.Server{
		Addr:    ":8080",
		Handler: routes,
	}

	log.Printf("starting HTTP server")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
