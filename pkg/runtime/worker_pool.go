package runtime

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fr0stylo/funcgo/pkg/utils"
)

type WorkerPool struct {
	workers     map[string]*Worker
	workerQueue chan string
	template    *WorkerOpts
	m           sync.RWMutex
	maxSize     int
}

type Runnable interface {
	Execute(any) ([]byte, error)
}

func NewWorkerPool(opts ...WorkerPoolOptsFn) *WorkerPool {
	cfg := defaultWorkerPoolOptions
	for _, opt := range opts {
		opt(cfg)
	}

	pool := &WorkerPool{
		template:    cfg.WorkerTemplate,
		workers:     map[string]*Worker{},
		maxSize:     cfg.MaxPoolSize,
		workerQueue: make(chan string, cfg.MaxPoolSize),
	}

	go pool.deregister()
	return pool
}

func (r *WorkerPool) Size() int {
	return len(r.workers)
}

func (r *WorkerPool) deregister() {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for range t.C {
		for n, w := range r.workers {
			if w.SinceLastExecution() > 200*time.Second && !w.IsBusy() {
				w.Stop()
				r.remove(n)
			}
		}
	}
}

func (r *WorkerPool) GetAvailable() *Worker {
	if len(r.workerQueue) > 0 {
		w := r.workers[<-r.workerQueue]
		if w == nil {
			return r.GetAvailable()
		}
		return w
	}

	if len(r.workers) <= r.maxSize {
		w := r.Push()

		return w
	}

	return r.GetAvailable()
}

func (r *WorkerPool) ExecOnWorker(obj any) ([]byte, error) {
	w := r.GetAvailable()
	defer func() {
		r.workerQueue <- w.name
	}()

	return w.Execute(obj)
}

func (r *WorkerPool) remove(name string) error {
	r.m.Lock()
	defer r.m.Unlock()

	delete(r.workers, name)
	log.Infof("[pool]: Removed %s workers left %v", name, r.workers)
	return nil
}

func (r *WorkerPool) Push() *Worker {
	name := fmt.Sprintf("%s-%s", utils.RandomString(8), utils.RandomString(8))
	log.Infof("[pool]: Started %s worker, current count: %v", name, len(r.workers))

	w := NewWorker(name, defaultIPManager.Acquire(), r.template)
	go w.Start(context.Background())

	r.m.Lock()
	defer r.m.Unlock()
	r.workers[name] = w

	return w
}
