package main

import "fmt"
import "os"

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("./devops --migrate")
		return
	}

	if args[1] == "--migrate" {
	}
}
