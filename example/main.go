package main

import (
	"fmt"
	"log"

	"github.com/pilu/treenotify"
)

func main() {
	w, err := treenotify.New()
	if err != nil {
		log.Fatal(err)
	}

	events, err := w.Watch(".")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("watching...\n")

	for {
		select {
		case event := <-events:
			fmt.Printf("%+v\n", event)
		}
	}
}
