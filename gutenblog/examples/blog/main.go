package main

import (
	"os"

	"katebuchepattison.com/gutenblog"
)

func main() {
	blog := gutenblog.New(
		"templates/layout.html.tmpl",
		"templates/index.html.tmpl",
		"templates/posts",
		"www", // For a simple website, OutDir and WebRoot will likely the same.
		"www", // If your blog is nested within an existing website, you must specify the root directory here.
	)

	args := os.Args
	if len(args) < 2 {
		panic("Not enough arguments\nUsage: build|serve")
	}

	switch args[1] {
	case "build":
		if err := blog.Generate(); err != nil {
			panic(err)
		}
	case "serve":
		blog.Serve("8080")
	default:
		panic("Unrecognized command\nUsage: build|serve")
	}
}
