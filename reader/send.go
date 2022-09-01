package reader

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jasontconnell/filesync/data"
)

func Send(clients []data.Client, duration time.Duration, files chan data.SyncFile) {
	go func() {
		for {
			// select {
			// case <-time.After(duration):
			select {
			case f := <-files:
				buf := bytes.NewBuffer(nil)
				enc := json.NewEncoder(buf)
				enc.Encode(f)
				for _, c := range clients {
					_, err := http.Post("http://"+c.Url+"/receive", "application/json", buf)

					if err != nil {
						log.Println("error sending file", f.RelativePath, "to", c.Url)
						continue
					}
				}
			default:
			}
		}
		// }
	}()
}
