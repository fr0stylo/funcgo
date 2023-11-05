package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

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
			// MaxConcurrency: 2,
			MaxConcurrency: 10,
			MainExec:       "/etc/function",
			RootFS:         "./fs",
			Files: runtime.FileList(
				// runtime.Files{From: "./wrapper.sh", To: "/etc/wrapper.sh"},
				runtime.Files{From: "./bin/function", To: "/etc/function"},
			),
		})

		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Hit me\n")
			r, _ := reader.ReadString('\n')

			fmt.Printf("%s\n", r)
			if r == "exit" {
				break
			}

			go mngr.Execute()
		}
	}
}

func must(name string, err error) {
	if err != nil {
		log.Fatal(name, err)
	}
}
