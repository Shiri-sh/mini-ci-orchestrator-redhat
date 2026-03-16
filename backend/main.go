package main

import (
	"fmt"
	"os"
	"path/filepath"
	"log"
	"net/http"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewApp(client *kubernetes.Clientset) *App {
	return &App{
		K8sClient: client,
		Builds: []Build{},
		NextBuildID: 1,
	}
}
func GetK8sClient() (*kubernetes.Clientset, error) {
	
	fmt.Println("Creating Kubernetes client...")
	config, err:= rest.InClusterConfig()
	if err != nil {
		home, _:= os.UserHomeDir()
		kubeconfig := filepath.Join(home, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err!=nil{
			return nil, fmt.Errorf("error couldn't find Kubernetes client config: %v \n", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes clientset: %v \n", err)
	}
	return clientset, nil
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
