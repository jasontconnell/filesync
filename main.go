package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/jasontconnell/filesync/conf"
	"github.com/jasontconnell/filesync/data"
	"github.com/jasontconnell/filesync/reader"
	"github.com/jasontconnell/filesync/writer"
)

func main() {
	fn := flag.String("c", "config.json", "the config file")
	flag.Parse()

	cfg, err := conf.LoadConfig(*fn)

	if !filepath.IsAbs(cfg.Path) {
		cfg.Path, _ = filepath.Abs(cfg.Path)
	}

	if err != nil {
		log.Println(err)
		return
	}

	sched, err := time.ParseDuration(cfg.Schedule)
	if err != nil {
		sched = time.Second * 60
		log.Println("parsing duration", err, sched)
	}

	if cfg.Role == "reader" {
		clients := []data.Client{}
		for _, c := range cfg.Clients {
			clients = append(clients, data.Client{Url: c})
		}

		done := make(chan bool)
		files := make(chan data.SyncFile)
		retry := make(chan data.SyncFile)
		reader.WatchDumb(cfg.Path, sched, cfg.Ignore, files, retry)
		reader.Send(clients, files)

		t := time.NewTicker(2 * time.Second)
		go func() {
			for tick := range t.C {
				fmt.Printf("\r%vClients: %d. File queue: %d. Retry queue: %d.\t\t",
					tick.Format("15:04:05"), len(cfg.Clients), len(files), len(retry))
			}
		}()

		<-done
	} else {
		h := writer.GetHandler(cfg.Path)
		log.Println("listening on", cfg.Bind, "writing to", cfg.Path)
		http.ListenAndServe(cfg.Bind, h)
	}
}
