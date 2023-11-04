package runtime

import "log"

type ManagerOpts struct {
	MaxConcurrency int
}

type Manager struct {
	pool           *WorkerPool
	maxConcurrency int
}

func NewManager(opts *ManagerOpts) *Manager {
	return &Manager{
		pool: NewWorkerPool(&WorkerPoolOpts{WorkerTemplate: &WorkerOpts{InitPath: "/etc/wrapper.sh", FilesToCopy: []Files{
			{From: "./wrapper.sh", To: "/etc/wrapper.sh"},
		}}}),
		maxConcurrency: opts.MaxConcurrency,
	}
}

func (r *Manager) Execute() {
	var w Runnable
	for w = r.pool.GetAvailable(); w == nil; w = r.pool.GetAvailable() {
		if r.maxConcurrency > r.pool.Size() {
			r.pool.Push()
		}

	}
	if w != nil {
		log.Printf("[manager]: Found %v", w)
		w.Execute()
		return
	}
	// 1. Find empty worker
	// 2. If no empty worker create one
	// 2.5 Wait till worker available
	// 3. Send message to worker
}
