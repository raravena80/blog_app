package main

import (
	"log"

	"github.com/raravena80/blog_app/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
