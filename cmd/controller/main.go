package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/fr0stylo/funcgo/pkg/apigw"
	"github.com/fr0stylo/funcgo/pkg/runtime"
)

func main() {
	log.Print("Running", os.Args)
	c := ""
	if len(os.Args) > 1 {
		c = os.Args[1]
	}

	switch c {
	case "container":
		containerInit()
	default:
		mngr := runtime.NewFunction(&runtime.FunctionOpts{
			MaxConcurrency: 10,
			MainExec:       "/etc/function",
			RootFS:         "./fs",
			Files: runtime.FileList(
				runtime.Files{From: "./bin/function", To: "/etc/function"},
			),
		})

		mux := mux.NewRouter()
		mux.Handle("/{id}", &apigw.Handler{Runner: mngr})
		// for {
		// 	reader := bufio.NewReader(os.Stdin)
		// 	fmt.Print("Hit me\n")
		// 	r, _ := reader.ReadString('\n')

		// 	fmt.Printf("%s\n", r)
		// 	if r == "exit" {
		// 		break
		// 	}

		// 	go mngr.Execute(map[string]string{"Url": "/asdasda", "Body": "ok"})
		// }

		http.ListenAndServe("0.0.0.0:8000", mux)
	}
}

func must(name string, err error) {
	if err != nil {
		log.Fatal(name, err)
	}
}
