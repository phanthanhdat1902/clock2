package main

import "fmt"

var R int = 99

func Workers(input chan int) { // Point #4 	// Point #1
	for i := 0; i < R; i++ { // Point #1
		go func(j int) {
			for {
				<-input
				fmt.Println(i)
			}
		}(i)
	} // Point #3
	input <- 10
	for {

	}
}

func main() {
	input := make(chan int)
	Workers(input)
	for {

	}
}
