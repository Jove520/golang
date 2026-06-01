package main

import "fmt"

type JobTx struct {
	store      *JobStore
	pending    map[int]Job
	committed  bool
	rolledBack bool
}

func (s *JobStore) Begin() *JobTx {
	return &JobTx{
		store:   s,
		pending: make(map[int]Job),
	}
}

func (tx *JobTx) Save(job Job) error {
	if tx.committed {
		return fmt.Errorf("transaction already committed")
	}
	if tx.rolledBack {
		return fmt.Errorf("transaction already rolled back")
	}

	tx.pending[job.ID] = job
	return nil
}

func (tx *JobTx) Commit() error {
	if tx.rolledBack {
		return fmt.Errorf("cannot commit rolled back transaction")
	}
	if tx.committed {
		return fmt.Errorf("transaction already committed")
	}

	tx.store.mu.Lock()
	defer tx.store.mu.Unlock()

	for id, job := range tx.pending {
		tx.store.jobs[id] = job
	}

	tx.committed = true
	return nil
}

func (tx *JobTx) Rollback() error {
	if tx.committed {
		return nil
	}
	if tx.rolledBack {
		return nil
	}

	tx.pending = nil
	tx.rolledBack = true
	return nil
}