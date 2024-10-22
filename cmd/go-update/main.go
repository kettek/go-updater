package main

import (
	"fmt"
	"os"

	"github.com/kettek/go-updater"
)

func main() {
	// Get parameters as source URL/filename, target filename, and optional PID
	if len(os.Args) < 3 {
		fmt.Println("Usage: update <source> <target> [PID]")
		return
	}
	source := os.Args[1]
	target := os.Args[2]
	pid := ""
	if len(os.Args) > 3 {
		pid = os.Args[3]
	}
	if err := updater.Update(source, target, pid); err != nil {
		fmt.Println(err)
	}
}
