package main

import (
	"log"

	"github.com/kabukky/klaus/filenames"
	"github.com/kabukky/klaus/klausignore"
	"github.com/kabukky/klaus/watcher"
)

func main() {
	done := make(chan bool)
	watcher.WatchRecursively(filenames.WorkingDirectory, klausignore.Ignore)
	defer watcher.Watcher.Close()
	log.Println("Watching", filenames.WorkingDirectory, "for changes.")
	// Wait indefinitely
	<-done
}
