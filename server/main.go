package main

import (
	"fmt"
	"os"
	"url-shortener/app"
)

//TODO: make docs
//TODO: add tests
//TODO: add readme file
//TODO: check links compared to provided ones in task
//TODO: check on format of shortened ID(if it should contain https://)
// TODO: support multiple environments with docker compose (instabug task)
//TODO: add graceful shutdown
func main() {
	app, err := app.NewApp()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app.Run()

}
