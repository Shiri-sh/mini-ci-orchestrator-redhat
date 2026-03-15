package main

import(
	"time"
)

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