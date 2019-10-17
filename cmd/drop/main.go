package main

import (
	"context"
	"github.com/diego1q2w/drop-box-it/pkg/drop/adapter/cli"
	"github.com/diego1q2w/drop-box-it/pkg/drop/app"
	"github.com/diego1q2w/drop-box-it/pkg/drop/infra"
	"log"
	"net/http"
	"os"
)

func main() {
	url, ok := os.LookupEnv("BOX_URL")
	if !ok {
		log.Fatal("BOX_URL env variable is required")
	}

	if len(os.Args) == 0 {
		log.Fatal("One argument is required which is the directory path you wish to sync")
	}

	boxClient := infra.NewBoxClient(http.DefaultClient, url)
	fileFetcher := infra.NewFileFetcher()

	dropApp := app.NewDropper(fileFetcher, boxClient, os.Args[1])
	cliDrop := cli.NewDropCLI(dropApp)

	ctx := context.Background()
	err := cliDrop.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}
