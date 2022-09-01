package reader

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jasontconnell/filesync/data"
	"github.com/rjeczalik/notify"
)

func Watch(path string, files chan data.SyncFile) {
	ch := make(chan notify.EventInfo, 1000)

	recpath := path + "\\.\\..."
	err := notify.Watch(recpath, ch, notify.All)
	if err != nil {
		log.Fatalf("error creating watch %s", err.Error())
	}

	go func() {
		for {
			select {
			case event := <-ch:
				err := getFiles(path, event.Path(), event, files)
				if err != nil {
					log.Println("error occurred reading files", event.Path(), err)
				}
			default:
			}
		}
	}()
}

func getFiles(start, path string, event notify.EventInfo, files chan data.SyncFile) error {
	del := event.Event() == notify.Remove
	rel := strings.Replace(path, start, "", -1)
	file := data.SyncFile{RelativePath: rel, Delete: del, Type: "file"}

	stat, err := os.Stat(path)
	if (err != nil || os.IsNotExist(err)) && !del {
		return err
	} else if del {
		files <- file
		return nil
	}

	if stat.IsDir() {
		file.Type = "directory"
	} else {
		contents, err := read(path)
		if err != nil {
			return err
		}
		file.Contents = contents
	}
	files <- file

	return nil
}

func read(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("couldn't read file %s. %w", path, err)
	}
	b64contents := base64.StdEncoding.EncodeToString(contents)
	return b64contents, nil
}
