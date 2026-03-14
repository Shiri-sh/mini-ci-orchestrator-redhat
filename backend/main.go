package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("../frontend")))

	// API endpoints
	http.HandleFunc("/builds", GetAllBuildsHandler)
	http.HandleFunc("/build/create", CreateBuildHandler)

	err:= http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}