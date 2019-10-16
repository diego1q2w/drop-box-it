package http

import (
	"encoding/base64"
	"github.com/diego1q2w/drop-box-it/pkg/box/domain"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

func DeleteDocumentHandler(service boxer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path64 := chi.URLParam(r, "path")

		path, err := base64.StdEncoding.DecodeString(path64)
		if err != nil {
			http.Error(w, "can't decode path", http.StatusBadRequest)
			return
		}

		err = service.DeleteDocuments(r.Context(), domain.Path(string(path)))
		if err != nil {
			log.Printf("Error writing document: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
