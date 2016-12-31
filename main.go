package main

import (
	"log"

	"github.com/kabukky/klaus/filenames"
	"github.com/kabukky/klaus/watcher"
)

func main() {
	done := make(chan bool)
	watcher.WatchRecursively(filenames.WorkingDirectory)
	defer watcher.Watcher.Close()
	log.Println("Watching", filenames.WorkingDirectory, "for changes.")
	// Wait indefinitely
	<-done
}
