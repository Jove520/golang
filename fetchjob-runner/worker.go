package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Worker struct {
	store		*JobStore
	fetcher 	*Fetcher
	timeout 	time.Duration
}

func NewWorker(store *JobStore, fetcher *Fetcher, timeout time.Duration) *Worker {
	return &Worker{
		store: store,
		fetcher: fetcher,
		timeout: timeout,
	}
}

func (w *Worker) Run(jobID int) {
	defer func() {
		if r := recover(); r != nil {
			w.store.UpdateStatus(jobID, StatusFailed, fmt.Sprintf("panic: %v", r))
			log.Printf("job %d recovered from panic: %v", jobID, r)
		}
	}()

	start := time.Now()
	defer func() {
		log.Printf("job %d finished, cost = %s", jobID, time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()

	job, ok := w.store.Get(jobID)
	if !ok {
		log.Printf("job %d not found", jobID)
		return
	}

	w.store.UpdateStatus(job.ID, StatusRunning, "")

	if job.URL == "panic" {
		panic("manual panic for testing recover")
	}

	err := w.fetcher.FetchToFile(ctx, job.URL, job.Path)

	if err != nil {
		job.Status = StatusFailed
		job.Error = err.Error()
	} else {
		job.Status = StatusSuccess
		job.Error = ""
	}

	tx := w.store.Begin()
	defer tx.Rollback()

	if err := tx.Save(job); err != nil {
		log.Printf("job %d save failed: %v", job.ID, err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("job %d commit failed: %v", job.ID, err)
		return
	}
}

