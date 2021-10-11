package main

import (
	"github.com/anschwa/gutenblog"
)

func main() {
	// Generate each blog in its appropriate place:
	rostock := gutenblog.New(
		"blog/rostock/layout.html.tmpl",
		"blog/rostock/index.html.tmpl",
		"blog/rostock/posts",
		"www/blog/rostock",
	)

	if err := rostock.Generate(); err != nil {
		panic(err)
	}
}
