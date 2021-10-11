package main

import (
	"os"

	"github.com/anschwa/gutenblog"
)

func main() {
	// Generate each blog in its appropriate place:
	rostock := gutenblog.New(
		"blog/rostock/layout.html.tmpl",
		"blog/rostock/index.html.tmpl",
		"blog/rostock/posts",
		"www/blog/rostock",
		"www",
	)

	args := os.Args
	if len(args) < 2 {
		panic("Not enough arguments\nUsage: build|serve")
	}

	switch args[1] {
	case "build":
		if err := rostock.Generate(); err != nil {
			panic(err)
		}
	case "serve":
		rostock.Serve("8080")
	default:
		panic("Unrecognized command\nUsage: build|serve")
	}
}
