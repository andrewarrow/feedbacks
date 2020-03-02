package main

import "fmt"
import "os"
import "os/exec"
import "github.com/andrewarrow/feedbacks/models"
import "github.com/andrewarrow/feedbacks/util"
import "github.com/andrewarrow/feedbacks/persist"

func main() {
	if util.InitConfig() == false {
		print("no config")
		return
	}
	args := os.Args
	if len(args) == 1 {
		fmt.Println("./devops --migrate")
		fmt.Println("./devops --sample")
		return
	}

	if args[1] == "--migrate" {
		cmd := exec.Command("mysql", "-u", "root", "feedbacks")
		cmd.Stdin, _ = os.Open("../migrations/first.sql")
		err := cmd.Run()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	} else if args[1] == "--sample" {
		db := persist.Connection()
		user := models.User{}
		user.Email = "test@user.com"
		user.Flavor = "admin"
		user.Phrase = "the rain stays mainly in maine"
		err := models.InsertUser(db, &user)
		if err != "" {
			fmt.Printf("%v\n", err)
		}
	}
}
