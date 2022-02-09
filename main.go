package main

import (
	"os"
	"os/signal"

	"katebuchepattison.com/gutenblog"
)

func main() {
	// Generate each blog in its appropriate place:
	rostock := gutenblog.New(
		"blog/rostock/layout.html.tmpl",
		"blog/rostock/index.html.tmpl",
		"blog/rostock/posts",
		"docs/blog/rostock",
		"docs",
	)

	ecuador := gutenblog.New(
		"blog/ecuador/layout.html.tmpl",
		"blog/ecuador/index.html.tmpl",
		"blog/ecuador/posts",
		"docs/blog/ecuador",
		"docs",
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
		if err := ecuador.Generate(); err != nil {
			panic(err)
		}

	case "serve":
		go rostock.Serve("8080")
		go ecuador.Serve("8081")

		// Wait for ^C
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		os.Exit(1)

	default:
		panic("Unrecognized command\nUsage: build|serve")
	}
}
