package reader

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/jasontconnell/filesync/data"
)

func WatchDumb(path string, schedule time.Duration, ignore []string, files, retry chan data.SyncFile) {
	go func() {
		t := time.NewTicker(schedule)
		defer t.Stop()

		lastUpdate := time.Time{}
		for {
			select {
			case ct := <-t.C:
				readUpdates(lastUpdate, path, ignore, files, retry)
				lastUpdate = ct
				retryUpdates(path, files, retry)
			default:
			}
		}
	}()
}

func retryUpdates(path string, files, retry chan data.SyncFile) {
	for ret := range retry {
		fp := filepath.Join(path, ret.RelativePath)
		contents, err := read(fp)
		if err != nil {
			retry <- ret
			log.Println("on retry, got error", err, fp)
			continue
		}
		ret.Contents = contents
		files <- ret
	}
}

func readUpdates(lastUpdate time.Time, path string, ignore []string, files, retry chan data.SyncFile) {
	filepath.Walk(path, func(fp string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		rel := strings.Replace(fp, path, "", 1)
		file := data.SyncFile{RelativePath: rel, Delete: false, Type: "file"}
		if info.ModTime().Local().After(lastUpdate) {
			contents, err := read(fp)
			if err != nil {
				log.Println("error reading", err)
				return nil
			}
			file.Contents = contents
			files <- file
		}

		return nil
	})
}
