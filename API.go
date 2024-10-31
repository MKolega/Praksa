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
	Error string `json:"error"`
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
	router.HandleFunc("/api/players/{id:[0-9]+}", makeHTTPHandlefunc(s.handleGetPlayerByID))
	router.HandleFunc("/api/ponude", makeHTTPHandlefunc(s.handleCreatePonuda)).Methods("POST")
	router.HandleFunc("/api/ponude/{id:[0-9]+}", makeHTTPHandlefunc(s.handleGetPonuda))
	router.HandleFunc("/api/deposit/{id:[0-9]+}", makeHTTPHandlefunc(s.handleDeposit)).Methods("POST")
	router.HandleFunc("/api/uplata/{id:[0-9]+}", makeHTTPHandlefunc(s.handleUplata)).Methods("POST")
	log.Println("JSON API Server is running on port: ", s.listenAddr)

	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		log.Fatal(err)
	}

}

func WriteJSON(w http.ResponseWriter, status int, v any) error {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func makeHTTPHandlefunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			_ = WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func (s *APIServer) ligeRequest(w http.ResponseWriter, _ *http.Request) error {
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

func (s *APIServer) handleGetPlayers(w http.ResponseWriter, _ *http.Request) error {
	player, err := s.store.GetPlayers()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, player)

}

func (s *APIServer) handleGetPlayerByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	player, err := s.store.GetPlayerByID(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, player)
}
func (s *APIServer) handleGetPonuda(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	ponuda, err := s.store.GetPonuda(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, ponuda)
}

func (s *APIServer) handleCreatePonuda(w http.ResponseWriter, r *http.Request) error {
	createPonudaReq := new(CreatePonudaRequest)
	if err := json.NewDecoder(r.Body).Decode(createPonudaReq); err != nil {
		return err
	}
	ponuda := NewPonuda(createPonudaReq.Broj, createPonudaReq.ID, createPonudaReq.Naziv, createPonudaReq.Vrijeme)
	if err := s.store.CreatePonuda(ponuda); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusCreated, ponuda)
}

func (s *APIServer) handleUplata(w http.ResponseWriter, r *http.Request) error {
	playerID, err := getID(r)
	if err != nil {
		return err
	}

	uplataReq := new(CreateUplataRequest)
	if err := json.NewDecoder(r.Body).Decode(&uplataReq); err != nil {
		return err
	}

	if err := s.store.CreateUplata(playerID, uplataReq.Amount, uplataReq.OdigraniPar); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, uplataReq)

}

func (s *APIServer) handleDeposit(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	depositRequest := new(DepositRequest)

	if err := json.NewDecoder(r.Body).Decode(&depositRequest); err != nil {
		return err
	}
	if err := s.store.Deposit(id, depositRequest.Amount); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, depositRequest)
}

func getID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return 0, fmt.Errorf("invalid player id: %s", vars["id"])
	}
	return id, nil
}
