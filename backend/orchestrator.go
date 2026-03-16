package main

import (
	"fmt"
	"context"
	"time"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	
	
)

func (app *App) TriggerBuild(b Build) {
	fmt.Printf("Starting Kubernetes job for repo: %s, branch: %s\n", b.Repo, b.Branch)
    
	ctx := context.Background()
    app.EnssurePVSExists()
	job := CloneSecurityJob(b)

	//create the job in Kubernetes
	_, err := app.K8sClient.BatchV1().Jobs("default").Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error creating Kubernetes job: %v\n", err)
		err = app.UpdateBuildStatus(b.ID,"failed")
		if err != nil {
			fmt.Printf("Error updating build status: %v\n", err)
		}
		return
	}
}

func FakeCloneJob(b Build) *batchv1.Job {
    return &batchv1.Job{
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
								"/bin/sh",
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
}
func int32Ptr(i int32) *int32 { return &i }

func CloneSecurityJob(b Build) *batchv1.Job {
	timestamp := time.Now().Unix()
	artifactName := fmt.Sprintf("security-%d-%d.json", b.ID, timestamp)
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("konflux-build-%d", b.ID),
			Labels: map[string]string{
				"app": "mini-ci",
				"build-id": fmt.Sprintf("%d", b.ID),
			},
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: int32Ptr(500), //clean up job after 500 seconds
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "repo-storage",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "artifact-storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "mini-ci-artifacts",
								},
							},
						},
					},
					
					InitContainers: []corev1.Container{
						{
							Name:  "git-clone",
							Image: "alpine/git",
							Command: []string{
								"sh",
								"-c",
								fmt.Sprintf("echo 'Cloning repo: %s'; git clone --depth 1 %s /workspace; echo 'Build finished'", b.Repo, b.Repo),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "repo-storage",
									MountPath: "/workspace",
								},
							},
						},
					},

					Containers: []corev1.Container{
						{
							Name:  "security-scan",
							Image: "trufflesecurity/trufflehog",
							Command: []string{
							"sh",
							"-c",
							fmt.Sprintf("trufflehog filesystem /workspace --json --fail > /artifacts/%s", artifactName),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:"repo-storage",
									MountPath: "/workspace",
								},
								{
									Name:"artifact-storage",
									MountPath: "/artifacts",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}
}
func (app *App) EnssurePVSExists() {
	pvsClient:= app.K8sClient.CoreV1().PersistentVolumeClaims("default")
	_, err := pvsClient.Get(context.Background(), "mini-ci-artifacts", metav1.GetOptions{})
	if err == nil {
		return//pvs exists
	}
	newPvs:= &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mini-ci-artifacts",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("500Mi"),
				},
			},
		},
	}
	_, err = pvsClient.Create(context.Background(), newPvs, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Error creating PersistentVolumeClaim: %v\n", err)
	}
}