package main

import (
	"fmt"

	expose "example.com/app/__doctest_internal_expose/greet"
)

func main() {
	fmt.Println(expose.Hello())
}
