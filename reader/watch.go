package reader

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jasontconnell/filesync/data"
	"github.com/rjeczalik/notify"
)

func Watch(path string, schedule time.Duration, ignore []string, files, retry chan data.SyncFile) {

	igmap := make(map[string]bool)
	for _, s := range ignore {
		igmap[s] = true
	}

	go loopWatch(path, igmap, files, retry)
	go loopRetry(path, files, retry)
}

func loopWatch(path string, igmap map[string]bool, files, retry chan data.SyncFile) {
	ch := make(chan notify.EventInfo, 1000)
	recpath := path + "\\.\\..."
	err := notify.Watch(recpath, ch, notify.All)
	if err != nil {
		log.Fatalf("error creating watch %s", err.Error())
	}
	defer notify.Stop(ch)

	for {
		select {
		case event := <-ch:
			err := getFiles(path, event.Path(), event, igmap, files, retry)
			log.Println("event", path, event.Event())
			if err != nil {
				log.Println("error occurred reading files", event.Path(), err)
			}
		default:
		}
	}
}

func loopRetry(path string, files, retry chan data.SyncFile) {
	for {
		select {
		case ret := <-retry:
			log.Println("retrying", ret.RelativePath)
			contents, err := read(filepath.Join(path, ret.RelativePath))
			if err != nil {
				log.Println("requeueing retry", ret.RelativePath)
				time.Sleep(time.Second * 5)
				retry <- ret
				continue
			}
			ret.Contents = contents
			files <- ret
		default:
		}
	}
}

func getFiles(start, path string, event notify.EventInfo, igmap map[string]bool, files, retry chan data.SyncFile) error {
	del := event.Event() == notify.Remove
	rel := strings.Replace(path, start, "", -1)
	file := data.SyncFile{RelativePath: rel, Delete: del, Type: "file"}

	_, fn := filepath.Split(path)
	if _, ok := igmap[fn]; ok {
		return nil
	}

	ext := filepath.Ext(path)
	if _, ok := igmap[ext]; ok {
		return nil
	}

	log.Println("file event raised", path, event.Event())

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
			retry <- file
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
