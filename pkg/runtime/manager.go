package runtime

import "log"

type ManagerOpts struct {
	MaxConcurrency int
	MainExec       string
	Files          []Files
	WorkerOptions  *WorkerOpts
	RootFS         string
}

type Manager struct {
	pool           *WorkerPool
	maxConcurrency int
}

func NewManager(opts *ManagerOpts) *Manager {
	return &Manager{
		pool:           NewWorkerPool(&WorkerPoolOpts{WorkerTemplate: &WorkerOpts{InitPath: opts.MainExec, FilesToCopy: opts.Files}}),
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

func must(name string, err error) {
	if err != nil {
		log.Fatal(name, err)
	}
}
