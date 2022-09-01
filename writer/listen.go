package writer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/jasontconnell/filesync/data"
)

func Listen(path string, files chan data.SyncFile) error {
	w, err := fsnotify.NewWatcher()

	if err != nil {
		return fmt.Errorf("Couldn't create watcher on %s. %w", path, err)
	}
	w.Add(path)

	go func() {
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					err := getFiles(path, event.Name, files)
					if err != nil {
						log.Println("error occurred reading files", err)
					}
				}
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	return nil
}

func getFiles(start, path string, files chan data.SyncFile) error {
	stat, err := os.Stat(path)

	if err != nil || os.IsNotExist(err) {
		return err
	}

	var ferr error
	if stat.IsDir() {
		ferr = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			rel := strings.Replace(p, start, "", -1)
			contents, err := read(p)
			if err != nil {
				return err
			}
			files <- data.SyncFile{RelativePath: rel, Contents: string(contents)}

			return nil
		})
	} else {
		rel := strings.Replace(path, start, "", -1)
		contents, err := read(path)
		files <- data.SyncFile{RelativePath: rel, Contents: string(contents)}
		ferr = err
	}

	return ferr
}

func read(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("couldn't read file %s. %w", path, err)
	}
	return string(contents), nil
}
