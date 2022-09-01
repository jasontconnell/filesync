package writer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/jasontconnell/filesync/data"
)

type ReaderHandler struct {
	http.Handler
	BasePath string
}

func GetHandler(path string) http.Handler {
	os.MkdirAll(path, os.ModePerm)

	h := ReaderHandler{BasePath: path}

	m := mux.NewRouter()
	m.HandleFunc("/receive", h.Receive)

	h.Handler = m
	return h
}

func readJson(r io.Reader, obj interface{}) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(obj)
	return err
}

func sendError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func (h ReaderHandler) Receive(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var file data.SyncFile
		err := readJson(req.Body, &file)
		if err != nil {
			sendError(w, fmt.Errorf("couldn't read json %w", err))
		}

		writeFile(filepath.Join(h.BasePath, file.RelativePath), file.Contents)
	}
}

func writeFile(path, contents string) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return fmt.Errorf("couldn't create dir tree %s. %w", path, err)
	}
	return os.WriteFile(path, []byte(contents), os.ModePerm)
}
