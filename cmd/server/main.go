package main

import (
	"flag"
	"log"

	"github.com/hackallcode/db-homework/internal/app/server"
)

func main() {
	params := server.Params{}
	flag.Int64Var(&params.Port, "port", 5000, "web port")
	flag.StringVar(&params.Url, "url", "/api", "web url")
	flag.BoolVar(&params.Reset, "reset", false, "reset db")
	flag.Parse()

	err := server.StartApp(params)
	if err != nil {
		log.Println(err)
		return
	}

	err = server.StopApp()
	if err != nil {
		log.Println(err)
		return
	}
}
