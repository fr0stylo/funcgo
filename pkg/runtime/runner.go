package runtime

import (
	"fmt"
	"log"
)

type FunctionOpts struct {
	MaxConcurrency int
	MainExec       string
	Files          []Files
	WorkerOptions  *WorkerOpts
	RootFS         string
}

type Function struct {
	pool           *WorkerPool
	maxConcurrency int
}

func NewFunction(opts *FunctionOpts) *Function {
	return &Function{
		pool:           NewWorkerPool(&WorkerPoolOpts{WorkerTemplate: &WorkerOpts{InitPath: opts.MainExec, FilesToCopy: opts.Files}}),
		maxConcurrency: opts.MaxConcurrency,
	}
}

func (r *Function) Execute(obj any) (any, error) {
	var w Runnable
	for w = r.pool.GetAvailable(); w == nil; w = r.pool.GetAvailable() {
		if r.maxConcurrency > r.pool.Size() {
			r.pool.Push()
		}

	}
	if w != nil {
		log.Printf("[manager]: Found %v", w)
		return w.Execute(obj)
	}
	// 1. Find empty worker
	// 2. If no empty worker create one
	// 2.5 Wait till worker available
	// 3. Send message to worker
	return nil, fmt.Errorf("failed to get lambda instance")
}

func must(name string, err error) {
	if err != nil {
		log.Fatal(name, err)
	}
}
