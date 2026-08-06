package main

import (
	"fmt"

	shim "example.com/realmod/__wrong_shim/leaf"
)

func main() {
	fmt.Println(shim.Bridge())
}
