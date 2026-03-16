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
	"strconv"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"context"
)

func NewApp(client *kubernetes.Clientset) *App {
	return &App{
		K8sClient: client,
		Builds: []Build{},
		NextBuildID: 1,
	}
}
func (app *App) WatchJobs() {
	//watch only to jobs with label app=mini-ci to avoid watching all jobs in the cluster
	watch, err := app.K8sClient.BatchV1().Jobs("default").Watch(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=mini-ci",
	})
	if err != nil {
		fmt.Printf("Error starting watch: %v\n", err)
		return
	}

	fmt.Println("Watching Kubernetes Jobs for updates...")

	for event := range watch.ResultChan() {
		job, ok := event.Object.(*batchv1.Job)
		if !ok { continue }

		buildIDStr := job.Labels["build-id"]
		buildID, _ := strconv.Atoi(buildIDStr)

		if job.Status.Active > 0 {
			app.UpdateBuildStatus(buildID, "running")
			fmt.Printf("Build %d is currently running...\n", buildID)
			continue 
		}

	
		if job.Status.Succeeded > 0 {
			app.UpdateBuildStatus(buildID, "success")
			fmt.Printf("Build %d finished successfully\n", buildID)
		} else if job.Status.Failed > 0 {
			app.UpdateBuildStatus(buildID, "failed")
			fmt.Printf("Build %d failed\n", buildID)
		}
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
	go app.WatchJobs()
	// API endpoints
	http.HandleFunc("/builds", app.GetAllBuildsHandler)
	http.HandleFunc("/build/create", app.CreateBuildHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
