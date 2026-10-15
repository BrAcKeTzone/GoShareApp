//go:build !android
// +build !android

package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "goshare" {
		fmt.Println("Usage: goshare goshare")
		os.Exit(1)
	}

	// Desktop mode: shares current working directory
	StartServer(false)
}