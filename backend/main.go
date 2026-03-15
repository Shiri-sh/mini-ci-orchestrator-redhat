package main

import (
	"log"
	"net/http"
	"k8s.io/client-go/kubernetes"
)

func NewApp(client *kubernetes.Clientset) *App {
	return &App{
		K8sClient: client,
		Builds: []Build{},
		NextBuildID: 1,
	}
}
func main() {
	http.Handle("/", http.FileServer(http.Dir("../frontend")))

	client, err:= GetK8sClient()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}
	app := NewApp(client)
	// API endpoints
	http.HandleFunc("/builds", app.GetAllBuildsHandler)
	http.HandleFunc("/build/create", app.CreateBuildHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	log.Println("Server started on port 8080")

}