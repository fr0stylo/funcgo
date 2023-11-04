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
	mngr := runtime.NewManager(&runtime.ManagerOpts{
		// MaxConcurrency: 2,
		MaxConcurrency: 10,
	})

	switch c {
	case "container":
		containerInit()
	default:
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
