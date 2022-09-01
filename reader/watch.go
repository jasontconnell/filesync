package reader

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jasontconnell/filesync/data"
	"github.com/rjeczalik/notify"
)

func Watch(path string, files chan data.SyncFile) {
	ch := make(chan notify.EventInfo, 1000)

	recpath := path + "\\.\\..."
	err := notify.Watch(recpath, ch, notify.Write)
	if err != nil {
		log.Fatalf("error creating watch %s", err.Error())
	}

	go func() {
		for {
			select {
			case event := <-ch:
				err := getFiles(path, event.Path(), files)
				if err != nil {
					log.Println("error occurred reading files", err)
				}
			default:
			}
		}
	}()
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
