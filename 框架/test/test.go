package main

import (
	"fmt"
)

func a() (int, int) {
	return 1, 2
}
func b() (int, int) {
	return a()
}
func main() {
	fmt.Print(b())
}
