package main

type JobStatus string

const (
	StatusPending JobStatus = "pending"
	StatusRunning JobStatus = "running"
	StatusSuccess JobStatus = "success"
	StatusFailed JobStatus = "failed"
)

type Job struct {
	ID		int
	URL		string
	Path	string
	Status	JobStatus
	Error	string
}
