package main

import (
	"fmt"
	"errors"
)

var (
	validationError = errors.New("validation error")
	errorNotFound   = errors.New("not found error")
)

func main() {
	fmt.Println("Hello, World!")
}