package main

import (
	"github.com/MKolega/Praksa/external"
	"github.com/MKolega/Praksa/internal/API"
	"github.com/MKolega/Praksa/internal/storage"
	"log"
)

func main() {
	store, err := storage.NewPostGresStore()
	get, err := external.NewPostGresGet()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	server := API.NewApiServer(":8080", store, get)
	server.Run()

}
