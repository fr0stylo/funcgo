package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/fr0stylo/funcgo/pkg/apigw"
	"github.com/fr0stylo/funcgo/pkg/runtime"
)

var (
	z, _ = zap.NewProduction()
	log  = z.Sugar()
)

func main() {
	log.Info("Running", os.Args)
	c := ""
	if len(os.Args) > 1 {
		c = os.Args[1]
	}

	host := runtime.NewHost()
	host.InsertFunction("function", runtime.NewFunction(&runtime.FunctionOpts{
		MaxConcurrency: 20,
		MainExec:       "/etc/function",
		RootFS:         "./fs",
		Files: runtime.FileList(
			runtime.Files{From: "./bin/function", To: "/etc/function"},
		),
	}))

	host.InsertFunction("function2", runtime.NewFunction(&runtime.FunctionOpts{
		MaxConcurrency: 20,
		MainExec:       "/etc/function",
		RootFS:         "./fs",
		Files: runtime.FileList(
			runtime.Files{From: "./bin/function2", To: "/etc/function"},
		),
	}))

	switch c {
	case "container":
		containerInit()
	default:
		mux := mux.NewRouter()
		mux.Handle("/{id}", &apigw.Handler{Runner: host, FunctionName: "function"}).Methods(http.MethodGet)
		mux.Handle("/{id}", &apigw.Handler{Runner: host, FunctionName: "function2"}).Methods(http.MethodPost)

		log.Info("Listening on :8000")
		http.ListenAndServe("0.0.0.0:8000", mux)
	}
}

func must(name string, err error) {
	if err != nil {
		log.Fatal(name, err)
	}
}
