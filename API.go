package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

type APIError struct {
	Error string
}
type apiFunc func(http.ResponseWriter, *http.Request) error

func NewApiServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/api/lige", makeHTTPHandlefunc(s.ligeRequest))
	router.HandleFunc("/api/players", makeHTTPHandlefunc(s.handlePlayer))
	router.HandleFunc("/api/ponude/{id:[0-9]+}", makeHTTPHandlefunc(s.handleGetPonuda))
	log.Println("JSON API Server is running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)

}

func WriteJSON(w http.ResponseWriter, status int, v any) error {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func makeHTTPHandlefunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})

		}
	}
}

func (s *APIServer) ligeRequest(w http.ResponseWriter, r *http.Request) error {
	lige, err := s.store.GetLige()

	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, lige)
}

func (s *APIServer) handlePlayer(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetPlayers(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreatePlayer(w, r)

	}
	return fmt.Errorf("method not allowed %s", r.Method)
}
func (s *APIServer) handleCreatePlayer(w http.ResponseWriter, r *http.Request) error {
	createPlayerReq := new(CreatePlayerRequest)
	if err := json.NewDecoder(r.Body).Decode(createPlayerReq); err != nil {
		return err
	}
	player := NewPlayer(createPlayerReq.Username, createPlayerReq.Password)
	if err := s.store.CreatePlayer(player); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusCreated, player)
}

func (s *APIServer) handleGetPlayers(w http.ResponseWriter, r *http.Request) error {
	player, err := s.store.GetPlayers()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, player)

}

func (s *APIServer) handleGetPonuda(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return fmt.Errorf("invalid ponuda: %s", vars["id"])
	}
	ponuda, err := s.store.GetPonuda(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, ponuda)
}
