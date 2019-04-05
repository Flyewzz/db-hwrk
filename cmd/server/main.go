package main

import (
	"log"
	"os"

	"github.com/hackallcode/db-homework/internal/app/server"
)

func main() {
	params := server.Params{
		Port:      os.Getenv("PORT"),
		ApiPrefix: "/api",
	}
	if params.Port == "" {
		params.Port = "5000"
	}

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
