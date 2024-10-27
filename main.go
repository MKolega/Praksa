package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Lige struct {
	Lige []struct {
		Naziv   string `json:"naziv"`
		Razrade []struct {
			Tipovi []struct {
				Naziv string `json:"naziv"`
			} `json:"tipovi"`
			Ponude []int `json:"ponude"`
		} `json:"razrade"`
	} `json:"lige"`
}

type Ponude []struct {
	Broj     string    `json:"broj"`
	ID       int       `json:"id"`
	Naziv    string    `json:"naziv"`
	Vrijeme  time.Time `json:"vrijeme"`
	Tecajevi []struct {
		Tecaj float64 `json:"tecaj"`
		Naziv string  `json:"naziv"`
	} `json:"tecajevi"`
	TvKanal       string `json:"tv_kanal,omitempty"`
	ImaStatistiku bool   `json:"ima_statistiku,omitempty"`
}

func main() {
	urls := []string{
		"https://minus5-dev-test.s3.eu-central-1.amazonaws.com/lige.json",
		"https://minus5-dev-test.s3.eu-central-1.amazonaws.com/ponude.json",
	}
	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			if url == "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/lige.json" {
				var lige Lige
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				if err := json.Unmarshal(body, &lige); err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Lige: %v\n", lige)

			} else if url == "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/ponude.json" {
				var ponude Ponude
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				if err := json.Unmarshal(body, &ponude); err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Ponude: %v\n", ponude)
			} else {
				log.Fatalf("unknown URL: %s", url)
			}

		} else {
			log.Fatalf("failed to fetch data: %s", resp.Status)
		}
	}
}
