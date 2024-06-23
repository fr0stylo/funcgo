package runtime

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
		pool: NewWorkerPool(
			WithWorkerTemplate(&WorkerOpts{InitPath: opts.MainExec, FilesToCopy: opts.Files}),
			WithMaxPoolSize(opts.MaxConcurrency)),
	}
}

func (r *Function) Execute(obj any) ([]byte, error) {
	return r.pool.ExecOnWorker(obj)
}

func must(name string, err error) {
	if err != nil {
		log.Fatal(name, err)
	}
}
