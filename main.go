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

func fetchData() (Lige, Ponude) {
	urls := []string{
		"https://minus5-dev-test.s3.eu-central-1.amazonaws.com/lige.json",
		"https://minus5-dev-test.s3.eu-central-1.amazonaws.com/ponude.json",
	}
	var lige Lige
	var ponude Ponude

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			if url == "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/lige.json" {

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				if err := json.Unmarshal(body, &lige); err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Lige: %v\n", lige)

			} else if url == "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/ponude.json" {
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
	return lige, ponude
}

func ligeRequest(lige Lige) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, liga := range lige.Lige {
			fmt.Fprintf(w, "Liga: %s\n", liga.Naziv)
			for _, razrada := range liga.Razrade {
				for _, ponuda := range razrada.Ponude {
					fmt.Fprintf(w, "Ponuda: %d\n", ponuda)
				}
			}
		}
	}
}

func ponudeRequest(ponude Ponude) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		for _, ponuda := range ponude {
			if fmt.Sprintf("/api/ponude/%d", ponuda.ID) == r.URL.Path {
				fmt.Fprintf(w, "Ponuda: %s\n", ponuda.Naziv)
			}

		}
	}
}

func addRequest(ponude Ponude, lige Lige) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/api/ponude/add" {
			var newPonuda struct {
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
			if err := json.NewDecoder(r.Body).Decode(&newPonuda); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ponude = append(ponude, newPonuda)
			w.WriteHeader(http.StatusCreated)
		} else if r.Method == "POST" && r.URL.Path == "/api/lige/add" {
			var newLiga struct {
				Naziv   string `json:"naziv"`
				Razrade []struct {
					Tipovi []struct {
						Naziv string `json:"naziv"`
					} `json:"tipovi"`
					Ponude []int `json:"ponude"`
				} `json:"razrade"`
			}
			if err := json.NewDecoder(r.Body).Decode(&newLiga); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			lige.Lige = append(lige.Lige, newLiga)
			w.WriteHeader(http.StatusCreated)
		} else {
			http.Error(w, "Invalid request method or path", http.StatusMethodNotAllowed)
		}
	}
}

func main() {
	lige, ponude := fetchData()
	http.HandleFunc("/api/lige", ligeRequest(lige))
	http.HandleFunc("/api/ponude/", ponudeRequest(ponude))
	http.HandleFunc("/api/ponude/add", addRequest(ponude, lige))
	http.HandleFunc("/api/lige/add", addRequest(ponude, lige))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
