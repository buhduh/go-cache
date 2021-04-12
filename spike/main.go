package main

import (
	"fmt"
)

type fooError struct{}

func (f fooError) Error() string {
	return "foo"
}

func isFooError(err error) bool {
	_, ok := err.(fooError)
	return ok
}

func main() {
	var err error = nil
	if isFooError(err) {
		fmt.Println("isFooError")
	} else {
		fmt.Println("not isFooError")
	}
}
