package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/komalkantmillan/urlshortener/handler"
	"github.com/komalkantmillan/urlshortener/service"
)

func main() {
	// Initialize the URL shortener service
	shortenerService := service.NewURLShortenerService()

	// Create a new router
	router := mux.NewRouter()

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	
	// Serve index.html for the root path
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// Register API handlers
	handler.RegisterHandlers(router, shortenerService)

	// Start the server
	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}