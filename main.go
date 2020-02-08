package main

import "github.com/andrewarrow/feedbacks/server"
import "time"
import "math/rand"
import "github.com/andrewarrow/feedbacks/util"
import "fmt"

func main() {
	rand.Seed(time.Now().UnixNano())
	if util.InitConfig() == false {
		print("no config")
		return
	}
	fmt.Println(util.AllConfig)
	server.Serve()
}
