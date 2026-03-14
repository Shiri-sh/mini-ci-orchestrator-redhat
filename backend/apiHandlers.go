package main

import(
	"encoding/json"
	"net/http"
	"time"
)

// API Handlers

// Get all builds
func GetAllBuildsHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetAllBuilds())
}

// Create a new build
func CreateBuildHandler(w http.ResponseWriter, r *http.Request) {
	var b Build

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	mu.Lock()

	b.ID = nextBuildID
	nextBuildID++
	b.Status = "pending"
	b.Repo = r.FormValue("repo")
	b.Branch = r.FormValue("branch")
	b.CreatedAt = time.Now()
	newBuild := AddBuild(b)
	
	mu.Unlock()
    
	go TriggerBuild(newBuild)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBuild)
}