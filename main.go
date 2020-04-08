package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/andrewarrow/feedbacks/server"
	"github.com/andrewarrow/feedbacks/util"

	e "github.com/andrewarrow/feedbacks/email"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	if util.InitConfig() == false {
		print("no config")
		return
	}
	if len(os.Args) == 4 && os.Args[1] == "email" {
		e.Send(os.Args[2], "andrew@many.pw", "welcome to socialdistance.app", e.MakeEmailHTML(string(os.Args[3])))
		return
	}
	fmt.Println(util.AllConfig)
	server.Serve()
}
