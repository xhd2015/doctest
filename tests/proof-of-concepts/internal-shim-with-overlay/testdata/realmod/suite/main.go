package main

import (
	"fmt"

	shim "example.com/realmod/http/__doctest_internal_shim/leaf"
)

func main() {
	fmt.Println(shim.Bridge())
}
