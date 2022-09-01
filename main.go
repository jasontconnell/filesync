package main

import (
	"flag"
	"log"
	"net/http"
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

	if err != nil {
		flag.PrintDefaults()
		return
	}

	if cfg.Role == "reader" {
		sched, _ := time.ParseDuration(cfg.Schedule)
		log.Println(sched)

		clients := []data.Client{}
		for _, c := range cfg.Clients {
			clients = append(clients, data.Client{Url: c})
		}

		done := make(chan bool)
		files := make(chan data.SyncFile)
		reader.Watch(cfg.Path, files)
		reader.Send(clients, files)
		<-done
	} else {
		h := writer.GetHandler(cfg.Path)
		log.Println("listening on", cfg.Bind, "writing to", cfg.Path)
		http.ListenAndServe(cfg.Bind, h)
	}
}
