// This sample program demonstrates how to create goroutines and

// how the scheduler behaves.

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// main is the entry point for all Go programs.
var wg sync.WaitGroup

func createPizza(pizza int) {
	defer wg.Done()
	time.Sleep(time.Second)
	fmt.Printf("Create %d\n", pizza)
}
func timeTrack(start time.Time, funName string) {
	elapsed := time.Since(start)
	fmt.Println(funName, "took", elapsed)
}
func main() {
	defer timeTrack(time.Now(), "Build Pizzas")
	runtime.GOMAXPROCS(3)
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go createPizza(i)
	}
	wg.Wait()
}
