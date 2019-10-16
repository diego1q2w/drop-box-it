package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/diego1q2w/drop-box-it/pkg/box/domain"
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//go:generate moq -out boxer_mock_test.go . boxer
type boxer interface {
	WriteDocuments(ctx context.Context, file domain.File) error
	DeleteDocuments(ctx context.Context, path domain.Path) error
}

type createFile struct {
	Content []byte `json:"content"`
	Mode    uint32 `json:"mode"`
}

func WriteDocumentHandler(service boxer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path64 := chi.URLParam(r, "path")

		path, err := base64.StdEncoding.DecodeString(path64)
		if err != nil {
			http.Error(w, "can't decode path", http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		var createFile createFile
		if err := json.Unmarshal(body, &createFile); err != nil {
			http.Error(w, "can't unmarshal body", http.StatusBadRequest)
			return
		}

		file := domain.File{
			Path:    domain.Path(path),
			Mode:    os.FileMode(createFile.Mode),
			Content: createFile.Content,
		}
		err = service.WriteDocuments(r.Context(), file)
		if err != nil {
			log.Printf("Error writing document: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
