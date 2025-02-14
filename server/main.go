package main

import (
	"fmt"
	"os"
	"url-shortener/app"
)

// TODO: support multiple environments with docker compose
// TODO: add graceful shutdown
func main() {
	app, err := app.NewApp()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app.Run()

}
