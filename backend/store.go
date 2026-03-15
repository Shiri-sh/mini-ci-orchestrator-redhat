package main

import (
	"k8s.io/client-go/kubernetes"
)

func NewApp(client *kubernetes.Clientset) *App {
	return &App{
		K8sClient: client,
		Builds: []Build{},
		NextBuildID: 1,
	}
}

func (app *App) AddBuild(b Build) Build {
	app.Mu.Lock()
	defer app.Mu.Unlock()
	b.ID = app.NextBuildID
	app.NextBuildID++
	app.Builds = append(app.Builds, b)
	return b
}

func (app *App) GetAllBuilds() []Build {
	app.Mu.Lock()
	defer app.Mu.Unlock()
	//returns a copy of the builds slice to prevent external modification
	copiedBuilds := make([]Build, len(app.Builds))
	copy(copiedBuilds, app.Builds)
	return copiedBuilds
}

func (app *App) UpdateBuildStatus(buildID int, status string) {
	app.Mu.Lock()
	defer app.Mu.Unlock()
	for i, b := range app.Builds {
		if b.ID == buildID {
			app.Builds[i].Status = status
			break
		}
	}
}