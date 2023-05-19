package reader

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jasontconnell/filesync/data"
)

func Send(clients []data.Client, files chan data.SyncFile) {
	go func() {
		for {
			select {
			case f := <-files:
				for _, c := range clients {
					buf := bytes.NewBuffer(nil)
					enc := json.NewEncoder(buf)
					enc.Encode(f)

					client := http.Client{}
					req, err := http.NewRequest(http.MethodPost, "http://"+c.Url+"/receive", buf)
					client.Timeout = time.Second * 10

					resp, err := client.Do(req)

					if err != nil {
						log.Println("error sending file", f.RelativePath, "to", c.Url)
						continue
					}

					if resp != nil && resp.StatusCode != 200 {
						log.Println("error sending file", f.RelativePath, "to", c.Url, "status code:", resp.StatusCode)
						continue
					}
				}
			default:
			}
		}
	}()
}
