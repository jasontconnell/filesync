package writer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/jasontconnell/filesync/data"
)

type WriterHandler struct {
	http.Handler
	BasePath string
}

func GetHandler(path string) http.Handler {
	os.MkdirAll(path, os.ModePerm)

	h := WriterHandler{BasePath: path}

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

func (h WriterHandler) Receive(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var file data.SyncFile
		err := readJson(req.Body, &file)
		if err != nil {
			sendError(w, fmt.Errorf("couldn't read json %w", err))
			return
		}

		err = writeFile(filepath.Join(h.BasePath, file.RelativePath), file.Contents, file.Type, file.Delete)
		if err != nil {
			sendError(w, fmt.Errorf("couldn't write file %s. %w", file.RelativePath, err))
			return
		}
	} else {
		w.Write([]byte("post only. see <a href=\"https://github.com/jasontconnell/filesync\" target=\"_blank\">https://github.com/jasontconnell/filesync</a>"))
	}
}

func writeFile(path, contents, ftype string, del bool) error {
	if del {
		return os.RemoveAll(path)
	}

	if ftype == "directory" {
		return os.MkdirAll(path, os.ModePerm)
	}

	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return fmt.Errorf("couldn't create dir tree %s. %w", path, err)
	}
	data, err := base64.StdEncoding.DecodeString(contents)
	if err != nil {
		return fmt.Errorf("couldn't decode base64 string %s. %w", contents, err)
	}
	return os.WriteFile(path, data, os.ModePerm)
}
