package main

import (
	// "time"
	"fmt"
	"context"
	"os"
	"path/filepath"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func (app *App) TriggerBuild(b Build) {
	fmt.Printf("Starting Kubernetes job for repo: %s, branch: %s\n", b.Repo, b.Branch)
    
	ctx := context.Background()

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("mini-ci-build-%d", b.ID),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "builder-container",
							Image: "alpine:latest",
							Command: []string{
								"bin/sh",
								"-c",
								fmt.Sprintf("echo 'Cloning repo: %s'; sleep 10; echo 'Build finished for repo: %s, branch: %s'", b.Repo, b.Repo, b.Branch),
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}
	errUpdateRunning := app.UpdateBuildStatus(b.ID,"Running")
	if errUpdateRunning != nil {
		fmt.Printf("Error updating build status: %v\n", errUpdateRunning)
	}

	//create the job in Kubernetes
	_, err := app.K8sClient.BatchV1().Jobs("default").Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error creating Kubernetes job: %v\n", err)
		err = app.UpdateBuildStatus(b.ID,"Failed")
		if err != nil {
			fmt.Printf("Error updating build status: %v\n", err)
		}
		return
	}
	errUpdateSuccess := app.UpdateBuildStatus(b.ID,"Succeeded")
	if errUpdateSuccess != nil {
		fmt.Printf("Error updating build status: %v\n", errUpdateSuccess)
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

