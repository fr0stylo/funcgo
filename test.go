package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/fr0stylo/funcgo/pkg/runtime"
)

var wg sync.WaitGroup

func main() {
	m := runtime.NewIPManager("192.0.0.%d/24")

	for i := 0; i < 2500; i = i + 1 {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()

			ip := m.Acquire()
			log.Print(ii, " ", ip)
			time.Sleep(time.Duration(rand.Int()%100) * time.Millisecond)
			m.Release(ip)
		}(i)
	}

	wg.Wait()
}
