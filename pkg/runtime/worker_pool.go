package runtime

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fr0stylo/funcgo/pkg/utils"
)

type WorkerPool struct {
	workers  map[string]*Worker
	template *WorkerOpts
	m        sync.RWMutex
}

type WorkerPoolOpts struct {
	WorkerTemplate *WorkerOpts
}

type Runnable interface {
	Execute(any) (any, error)
}

func NewWorkerPool(opts *WorkerPoolOpts) *WorkerPool {
	pool := &WorkerPool{
		template: opts.WorkerTemplate,
		workers:  map[string]*Worker{},
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

func (r *WorkerPool) GetAvailable() Runnable {
	r.m.Lock()
	defer r.m.Unlock()

	for _, w := range r.workers {
		if !w.IsBusy() {
			return w
		}
	}

	return nil
}

func (r *WorkerPool) remove(name string) error {
	r.m.Lock()
	defer r.m.Unlock()

	delete(r.workers, name)
	log.Printf("[pool]: Removed %s workers left %v", name, r.workers)
	return nil
}

func (r *WorkerPool) Push() Runnable {
	r.m.Lock()
	name := fmt.Sprintf("%s-%s", utils.RandomString(8), utils.RandomString(8))

	w := NewWorker(name, defaultIPManager.Acquire(), r.template)
	w.Start(context.Background())

	r.workers[name] = w

	r.m.Unlock()
	return w
}
