package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "Hello Stderr!")
	fmt.Fprintln(os.Stderr, "Hello errors!")
}
