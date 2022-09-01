package writer

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/jasontconnell/filesync/data"
)

func Send(clients []data.Client, files chan data.SyncFile) {
	for {
		select {
		case f := <-files:
			buf := bytes.NewBuffer(nil)
			enc := json.NewEncoder(buf)
			enc.Encode(f)
			for _, c := range clients {
				http.Post(c.Url+"/receive", "application/json", buf)
			}

		}
	}
}
