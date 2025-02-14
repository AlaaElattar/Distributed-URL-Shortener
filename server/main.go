package main

import (
	"fmt"
	"os"
	"url-shortener/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app.Run()

}
