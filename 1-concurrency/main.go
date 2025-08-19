package main

import (
	"fmt"
	"math"
	"math/rand"
)

func randSlice(doubleCh chan int) {
	numCh := make(chan int, 10)
	mySlice := make([]int, 10)
	for idx := range mySlice {
		mySlice[idx] = rand.Intn(101)
	}
	go double(numCh, doubleCh)
	for _, value := range mySlice {
		numCh <- value
	}
}

func double(numCh chan int, doubleCh chan int) {
	for range 10 {
		doubleInt := int(math.Pow(float64(<-numCh), 2))
		doubleCh <- doubleInt
	}
}

func main() {
	doubleCh := make(chan int, 10)
	go randSlice(doubleCh)

	for range 10 {
		fmt.Println(<-doubleCh)
	}
}
