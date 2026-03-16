package main

import (
	"fmt"
	"context"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (app *App) TriggerBuild(b Build) {
	fmt.Printf("Starting Kubernetes job for repo: %s, branch: %s\n", b.Repo, b.Branch)
    
	ctx := context.Background()

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
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("konflux-build-%d", b.ID),
			Labels: map[string]string{
				"app": "mini-ci",
				"build-id": fmt.Sprintf("%d", b.ID),
			},
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: int32Ptr(100),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "repo-storage",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					
					InitContainers: []corev1.Container{
						{
							Name:  "git-clone",
							Image: "alpine/git:2.41.0",
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
							Image: "trufflesecurity/trufflehog:3.24.1",
							Args: []string{"filesystem","/workspace","--fail"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:"repo-storage",
									MountPath: "/workspace",
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
