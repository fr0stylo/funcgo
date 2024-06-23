package runtime

type WorkerPoolOpts struct {
	WorkerTemplate *WorkerOpts
	MaxPoolSize    int
}

var defaultWorkerPoolOptions = &WorkerPoolOpts{
	MaxPoolSize: 10,
	WorkerTemplate: &WorkerOpts{
		InitPath:    "",
		FilesToCopy: nil,
	},
}

type WorkerPoolOptsFn = func(w *WorkerPoolOpts)

func WithMaxPoolSize(i int) func(w *WorkerPoolOpts) {
	return func(w *WorkerPoolOpts) {
		w.MaxPoolSize = i
	}
}

func WithWorkerTemplate(i *WorkerOpts) func(w *WorkerPoolOpts) {
	return func(w *WorkerPoolOpts) {
		w.WorkerTemplate = i
	}
}
