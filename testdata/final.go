package main

import (
	"fmt"
)

type Result[T any] struct {
	Value T
	Err   error
}

func Divide(a, b int) Result[int] {
	if b == 0 {
		return Result[int]{Err: fmt.Errorf("divide by zero")}
	}
	return Result[int]{Value: a / b}
}

func main() {
	res := Divide(10, 2)
	if res.Err != nil {
		fmt.Println("Error:", res.Err)
	} else {
		fmt.Println("Result:", res.Value)
	}
}
