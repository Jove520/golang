package main

import "sync"

type JobStore struct {
	mu		sync.Mutex
	jobs	map[int]Job
}

func NewJobStore() *JobStore {
	return &JobStore{
		jobs: make(map[int]Job),
	}
}

func (s *JobStore) Add(job Job) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs[job.ID] = job
}

func (s *JobStore) UpdateStatus(id int, status JobStatus, errMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job := s.jobs[id]
	job.Status = status
	job.Error = errMsg
	s.jobs[id] = job
}

func (s *JobStore) Get(id int) (Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	return job, ok
}

func (s *JobStore) List() []Job {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		result = append(result, job)
	}

	return result
}