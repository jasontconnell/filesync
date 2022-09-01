package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/jasontconnell/filesync/conf"
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

	sched, err := time.ParseDuration(cfg.Schedule)
	fmt.Println(sched)

	done := make(chan bool)
	go writer.Listen(cfg.Path, nil)

	<-done
}
