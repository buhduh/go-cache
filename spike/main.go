package main

import "fmt"

func main() {
	foo := make(chan int, 1000)
	for i := 0; i < 1000; i++ {
		foo <- i
	}
loop:
	for {
		select {
		case tmp := <-foo:
			fmt.Printf("%d\n", tmp)
		default:
			fmt.Printf("default")
			break loop
		}
	}
}
