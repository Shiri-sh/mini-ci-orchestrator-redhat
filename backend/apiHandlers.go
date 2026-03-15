package main

import(
	"encoding/json"
	"net/http"
	"time"
	"fmt"
)

// API Handlers

// Get all builds
func GetAllBuildsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetAllBuilds())
}

// Create a new build
func CreateBuildHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating a new build")
    var br BuildRequest
	if err := json.NewDecoder(r.Body).Decode(&br); err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var b Build

	b.Status = "pending"
	b.Repo = br.Repo
	b.Branch = br.Branch
	b.CreatedAt = time.Now()
	newBuild := AddBuild(b)

	fmt.Printf("New build created: %+v\n", newBuild)
    
	go TriggerBuild(newBuild)
	
    fmt.Println("sending back response")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBuild)
}