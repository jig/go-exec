package main

import (
	"fmt"
	"os"
)

func main() {
	var s string
	_, _ = fmt.Scan(&s)
	if s != "example" {
		os.Exit(1)
	}
}
