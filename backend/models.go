package main

import(
	"time"
	"k8s.io/client-go/kubernetes"
	"sync"
)

type App struct {
	K8sClient *kubernetes.Clientset
	Builds []Build
	NextBuildID int
	Mu sync.Mutex
}

type Build struct {
	ID int `json:"id"`
	Repo string `json:"repo"`
	Branch string `json:"branch"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type BuildRequest struct {
	Repo string `json:"repo"`
	Branch string `json:"branch"`
}

