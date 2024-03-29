package main

import (
	"github.com/OneOfOne/xxhash"
	boxHttp "github.com/diego1q2w/drop-box-it/pkg/box/adapter/http"
	"github.com/diego1q2w/drop-box-it/pkg/box/app"
	"github.com/diego1q2w/drop-box-it/pkg/box/infra"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

const (
	rootDir      = "./destDir"
	numOfWorkers = 10
)

func main() {
	file := infra.NewFileBox()
	service := app.NewBox(file, xxHash, rootDir, numOfWorkers)
	writeHandler := boxHttp.WriteDocumentHandler(service)
	writeHandler = boxHttp.UncompressHandler(writeHandler)
	deleteHandler := boxHttp.DeleteDocumentHandler(service)

	mux := chi.NewMux()
	mux.Post("/document/{path}", writeHandler)
	mux.Delete("/document/{path}", deleteHandler)

	log.Println("Starting box server...")
	if err := http.ListenAndServe("0.0.0.0:80", mux); err != nil {
		log.Fatalf("unable to start the server: %s", err)
	}
}

func xxHash(s string) uint64 {
	return xxhash.Checksum64([]byte(s))
}
