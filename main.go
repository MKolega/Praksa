package main

import (
	"github.com/MKolega/Praksa/internal/API"
	"github.com/MKolega/Praksa/internal/storage"
	"log"
)

func main() {
	store, err := storage.NewPostGresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	server := API.NewApiServer(":8080", store)
	server.Run()

}
