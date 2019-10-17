package main

import (
	boxHttp "github.com/diego1q2w/drop-box-it/pkg/box/adapter/http"
	"github.com/diego1q2w/drop-box-it/pkg/box/app"
	"github.com/diego1q2w/drop-box-it/pkg/box/infra"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

const rootDir = "./destDir"

func main() {
	file := infra.NewFileBox()
	service := app.NewBox(file, rootDir)
	writeHandler := boxHttp.WriteDocumentHandler(service)
	deleteHandler := boxHttp.DeleteDocumentHandler(service)

	mux := chi.NewMux()
	mux.Post("/document/{path}", writeHandler)
	mux.Delete("/document/{path}", deleteHandler)

	log.Println("Starting box server...")
	if err := http.ListenAndServe("0.0.0.0:80", mux); err != nil {
		log.Fatalf("unable to start the server: %s", err)
	}
}
