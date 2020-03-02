package main

import "fmt"
import "os"
import "os/exec"

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("./devops --migrate")
		return
	}

	if args[1] == "--migrate" {
		cmd := exec.Command("mysql", "-u", "root", "feedbacks")
		cmd.Stdin, _ = os.Open("../migrations/first.sql")
		err := cmd.Run()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
}
