package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	if err := os.MkdirAll("output", 0755); err != nil {
		log.Fatalf("create output dir: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("hello from fetchjob runner\n"))
		case "/fail":
			http.Error(w, "server error", http.StatusInternalServerError)
		case "/slow":
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("slow response\n"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	store := NewJobStore()
	fetcher := NewFetcher()
	worker := NewWorker(store, fetcher, 1*time.Second)

	jobs := []Job{
		{
			ID:     1,
			URL:    server.URL + "/ok",
			Path:   "output/job-1.txt",
			Status: StatusPending,
		},
		{
			ID:     2,
			URL:    server.URL + "/fail",
			Path:   "output/job-2.txt",
			Status: StatusPending,
		},
		{
			ID:     3,
			URL:    server.URL + "/slow",
			Path:   "output/job-3.txt",
			Status: StatusPending,
		},
		{
			ID:     4,
			URL:    "panic",
			Path:   "output/job-4.txt",
			Status: StatusPending,
		},
	}

	for _, job := range jobs {
		store.Add(job)
	}

	for _, job := range jobs {
		worker.Run(job.ID)
	}

	fmt.Println("final job statuses:")
	for _, job := range store.List() {
		fmt.Printf("job=%d status=%s error=%q path=%s\n", job.ID, job.Status, job.Error, job.Path)
	}
}