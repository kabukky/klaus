package watcher

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"github.com/kabukky/klaus/builder"
	"github.com/kabukky/klaus/helper"
	"github.com/kabukky/klaus/klausignore"
	"github.com/kabukky/klaus/runner"
)

var (
	Watcher *fsnotify.Watcher
)

func init() {
	// Prepare watcher
	var err error
	Watcher, err = newWatcher()
	if err != nil {
		log.Fatalln("Could not create watcher:", err)
	}
}

func WatchRecursively(path string, ignoreDirs []string) error {
	// Watch root directory and all subdirectories in the given path
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			for _, d := range ignoreDirs {
				if d == info.Name() {
					return filepath.SkipDir
				}
			}

			err := Watcher.Add(filePath)
			if err != nil {
				return err
			}
		} else {
			// The directory we just added might have source files in it
			checkSourceFile(filePath)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func stopWatching(path string) error {
	return Watcher.Remove(path)
}

func checkSourceFile(path string) {
	if filepath.Ext(path) == ".go" {
		// A source file was modified. Rebuild.
		rebuild()
	}
}

func rebuild() {
	go func() {
		output, err := builder.Build()
		if err != nil {
			if err != builder.BuildAlreayUnderWayError {
				log.Printf("\033[91mError while building:\033[0m %s", output)
			}
		} else {
			log.Println("\033[92mBuild completed successfully.\033[0m")
			err := runner.Run()
			if err != nil {
				log.Printf("\033[91mCould not start binary:\033[0m %s", err.Error())
			}
		}
	}()
}

func newWatcher() (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					// A file was Modified
					checkSourceFile(event.Name)
				} else if event.Op&fsnotify.Create == fsnotify.Create {
					if helper.IsDirectory(event.Name) {
						// A new directory was created. Watch it.
						err := WatchRecursively(event.Name, klausignore.Ignore)
						if err != nil {
							log.Println("Error while while trying to add created directory to watcher:", err)
						}
					} else {
						checkSourceFile(event.Name)
					}
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					if helper.IsDirectory(event.Name) {
						// A directory was removed. Stop watching it.
						err := stopWatching(event.Name)
						if err != nil {
							log.Println("Error while while trying to remove directory from watcher:", err)
						}
					} else {
						checkSourceFile(event.Name)
					}
				}
			case err := <-watcher.Errors:
				log.Println("Error while watching directory:", err)
			}
		}
	}()
	return watcher, nil
}
