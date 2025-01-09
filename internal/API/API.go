package API

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MKolega/Praksa/internal/shared"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

type APIServer struct {
	listenAddr string
	store      shared.Storage
	get        shared.Get
}

type APIError struct {
	Error string `json:"error"`
}

func NewApiServer(listenAddr string, store shared.Storage, get shared.Get) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
		get:        get,
	}
}

func (s *APIServer) Run() {

	ligeURL := "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/lige.json"
	err := s.FetchAndInsertLigeDataToDB(ligeURL)
	if err != nil {
		log.Fatal("failed to insert lige data: ", err)
	}
	ponudeURL := "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/ponude.json"
	err = s.FetchAndInsertPonudeDataToDB(ponudeURL)
	if err != nil {
		log.Fatal("failed to insert ponude data: ", err)
	}

	router := mux.NewRouter()
	router.Use(enableCors)
	router.HandleFunc("/api/lige", makeHTTPHandlefunc(s.HandleGetLige))
	router.HandleFunc("/api/players", makeHTTPHandlefunc(s.handlePlayer))
	router.HandleFunc("/api/players/{id:[0-9]+}", makeHTTPHandlefunc(s.handleGetPlayerByID))
	router.HandleFunc("/api/login", makeHTTPHandlefunc(s.handleLogin))
	router.HandleFunc("/api/ponude", makeHTTPHandlefunc(s.handlePonude))
	router.HandleFunc("/api/ponude/{id:[0-9]+}", makeHTTPHandlefunc(s.handleGetPonuda))
	router.HandleFunc("/api/deposit/{id:[0-9]+}", makeHTTPHandlefunc(s.handleDeposit)).Methods("POST")
	router.HandleFunc("/api/uplata/{id:[0-9]+}", makeHTTPHandlefunc(s.handleUplata)).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./client/build")))

	log.Println("JSON API Server is running on port: ", s.listenAddr)

	err = http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		log.Fatal(err)
	}

}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandlefunc(handlerFunc func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := handlerFunc(w, r)

		if err != nil {

			switch e := err.(type) {
			case *shared.UserError:

				log.Printf("Bad request: %v", err)
				_ = WriteJSON(w, http.StatusBadRequest, APIError{Error: e.Message})
			case *shared.InternalError:

				log.Printf("Internal server error: %v", err)
				_ = WriteJSON(w, http.StatusInternalServerError, APIError{Error: e.Message})
			default:

				log.Printf("Unexpected error: %v", err)
				_ = WriteJSON(w, http.StatusInternalServerError, APIError{Error: "Unexpected error"})
			}
		}
	}
}

func (s *APIServer) FetchData(url string, decodeFunc func(io.Reader) error) error {
	r, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data from %s: %v", url, err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Println("failed to close response body:", err)
		}
	}(r.Body)

	if err := decodeFunc(r.Body); err != nil {
		return fmt.Errorf("failed to decode or process data from %s: %v", url, err)
	}

	log.Printf("Successfully processed data from %s.", url)
	return nil
}

func (s *APIServer) FetchAndInsertLigeDataToDB(url string) error {

	decodeFunc := func(body io.Reader) error {
		var jsonData shared.JsonData
		if err := json.NewDecoder(body).Decode(&jsonData); err != nil {
			return fmt.Errorf("failed to decode Lige JSON: %v", err)
		}

		for _, liga := range jsonData.Lige {
			ligaID, err := s.store.CreateLiga(liga.Naziv)
			if err != nil {
				log.Printf("failed to insert liga %s: %v", liga.Naziv, err)
				continue
			}

			for _, razrada := range liga.Razrade {
				razradaID, err := s.store.CreateRazrada(ligaID, razrada.Ponude)
				if err != nil {
					log.Printf("failed to insert razrada for liga %s: %v", liga.Naziv, err)
					continue
				}

				for _, tip := range razrada.Tipovi {
					err := s.store.CreateTipovi(razradaID, tip.Naziv)
					if err != nil {
						log.Printf("failed to insert tip %s for razrada %d: %v", tip.Naziv, razradaID, err)

					}
				}
			}
		}
		log.Println("Successfully updated Lige data.")

		return nil
	}
	return s.FetchData(url, decodeFunc)
}

func (s *APIServer) FetchAndInsertPonudeDataToDB(url string) error {
	decodeFunc := func(body io.Reader) error {
		var jsonData []shared.Ponude
		if err := json.NewDecoder(body).Decode(&jsonData); err != nil {
			return fmt.Errorf("failed to decode Ponude JSON: %v", err)
		}

		for _, ponuda := range jsonData {
			if err := s.store.CreatePonuda(&ponuda); err != nil {
				log.Printf("failed to insert ponuda with ID %d: %v", ponuda.ID, err)
				continue
			}

			for _, tecaj := range ponuda.Tecajevi {
				if err := s.store.CreateTecaj(ponuda.ID, tecaj.Tecaj, tecaj.Naziv); err != nil {
					log.Printf("failed to insert tecaj '%s' for ponuda ID %d: %v", tecaj.Naziv, ponuda.ID, err)
				}
			}
		}

		log.Println("Successfully updated Ponude data.")
		return nil
	}
	return s.FetchData(url, decodeFunc)
}

func (s *APIServer) HandleGetLige(w http.ResponseWriter, _ *http.Request) error {
	lige, err := s.get.GetLige()

	if err != nil {
		return &shared.InternalError{Message: fmt.Sprintf("failed to get lige: %v", err)}
	}
	return WriteJSON(w, http.StatusOK, lige)
}

func (s *APIServer) handlePlayer(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetPlayers(w, r)
	case "POST":
		return s.handleCreatePlayer(w, r)
	case "PUT":
		return s.handlePasswordReset(w, r)
	case "DELETE":
		return s.handleDeleteUser(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

func (s *APIServer) handlePonude(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handeGetAllPonude(w, r)
	case "POST":
		return s.handleCreatePonuda(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

func (s *APIServer) handleCreatePlayer(w http.ResponseWriter, r *http.Request) error {
	createPlayerReq := new(shared.CreatePlayerRequest)
	if err := json.NewDecoder(r.Body).Decode(createPlayerReq); err != nil {
		return &shared.UserError{Message: fmt.Sprintf("failed to decode player data: %v", err)}
	}
	player := shared.NewPlayer(createPlayerReq.Username, createPlayerReq.Password)
	if err := s.store.CreatePlayer(player); err != nil {
		return &shared.InternalError{Message: fmt.Sprintf("failed to create player: %v", err)}
	}
	return WriteJSON(w, http.StatusCreated, player)
}

func (s *APIServer) handleGetPlayers(w http.ResponseWriter, _ *http.Request) error {
	player, err := s.get.GetPlayers()
	if err != nil {
		return &shared.InternalError{Message: fmt.Sprintf("failed to get players: %v", err)}
	}

	return WriteJSON(w, http.StatusOK, player)

}

func (s *APIServer) handlePasswordReset(w http.ResponseWriter, r *http.Request) error {
	resetRequest := new(shared.CreatePlayerRequest)
	if err := json.NewDecoder(r.Body).Decode(resetRequest); err != nil {
		return &shared.UserError{Message: fmt.Sprintf("failed to decode player data: %v", err)}
	}
	if err := s.store.ResetPassword(resetRequest.Username, resetRequest.Password); err != nil {
		return &shared.InternalError{Message: fmt.Sprintf("failed to reset password: %v", err)}
	}
	return WriteJSON(w, http.StatusOK, resetRequest)
}

func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return &shared.UserError{Message: fmt.Sprintf("invalid player id: %v", err)}
	}

	if err := s.store.DeleteUser(id); err != nil {
		return &shared.InternalError{Message: fmt.Sprintf("failed to delete user: %v", err)}
	}
	return WriteJSON(w, http.StatusOK, id)

}
func (s *APIServer) handleGetPlayerByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return &shared.UserError{Message: fmt.Sprintf("invalid player id: %v", err)}
	}
	player, err := s.get.GetPlayerByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &shared.UserError{Message: fmt.Sprintf("Player with ID %d not found: %v", id, err)}
		}
		return &shared.InternalError{Message: fmt.Sprintf("failed to get player by id %d: %v", id, err)}
	}
	return WriteJSON(w, http.StatusOK, player)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	loginReq := new(shared.Player)
	if err := json.NewDecoder(r.Body).Decode(loginReq); err != nil {
		return &shared.UserError{Message: fmt.Sprintf("failed to decode login data: %v", err)}
	}

	player, err := s.get.GetLogin(loginReq.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &shared.UserError{Message: fmt.Sprintf("Player with username %s not found", loginReq.Username)}
		}

		return &shared.InternalError{Message: fmt.Sprintf("Login failed with username: %v", err)}
	}
	if player.Password == loginReq.Password {
		return WriteJSON(w, http.StatusOK, player)
	}

	return fmt.Errorf("invalid password")
}

func (s *APIServer) handleGetPonuda(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return &shared.UserError{Message: fmt.Sprintf("invalid ponuda id: %v", err)}
	}
	ponuda, err := s.get.GetPonuda(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &shared.UserError{Message: fmt.Sprintf("ponuda with id %d not found", id)}
		}
		return &shared.InternalError{Message: fmt.Sprintf("failed to get ponuda by id %d: %v", id, err)}
	}
	return WriteJSON(w, http.StatusOK, ponuda)
}

func (s *APIServer) handleCreatePonuda(w http.ResponseWriter, r *http.Request) error {
	createPonudaReq := new(shared.CreatePonudaRequest)
	if err := json.NewDecoder(r.Body).Decode(createPonudaReq); err != nil {
		return &shared.UserError{Message: fmt.Sprintf("failed to decode ponuda data: %v", err)}
	}

	ponuda := shared.NewPonuda(createPonudaReq.Broj, createPonudaReq.ID, createPonudaReq.Naziv, createPonudaReq.Vrijeme, createPonudaReq.TvKanal, createPonudaReq.ImaStatistiku)
	if err := s.store.CreatePonuda(ponuda); err != nil {
		return &shared.InternalError{Message: fmt.Sprintf("failed to create ponuda: %v", err)}
	}

	for _, tecaj := range createPonudaReq.Tecajevi {
		if err := s.store.CreateTecaj(createPonudaReq.ID, tecaj.Tecaj, tecaj.Naziv); err != nil {
			return &shared.InternalError{Message: fmt.Sprintf("failed to create tecaj: %v", err)}
		}
	}

	return WriteJSON(w, http.StatusCreated, createPonudaReq)
}
func (s *APIServer) handleUplata(w http.ResponseWriter, r *http.Request) error {
	playerID, err := getID(r)
	if err != nil {
		return &shared.UserError{Message: fmt.Sprintf("invalid player id: %v", err)}
	}

	uplataReq := new(shared.CreateUplataRequest)
	if err := json.NewDecoder(r.Body).Decode(&uplataReq); err != nil {
		return &shared.UserError{Message: fmt.Sprintf("failed to decode uplata data: %v", err)}
	}

	if err := s.store.CreateUplata(playerID, uplataReq.Amount, uplataReq.OdigraniPar); err != nil {
		return &shared.InternalError{Message: fmt.Sprintf("failed to create uplata: %v", err)}
	}

	return WriteJSON(w, http.StatusOK, uplataReq)

}

func (s *APIServer) handleDeposit(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return &shared.UserError{Message: fmt.Sprintf("invalid player id: %v", err)}
	}
	depositRequest := new(shared.DepositRequest)

	if err := json.NewDecoder(r.Body).Decode(&depositRequest); err != nil {
		return &shared.UserError{Message: fmt.Sprintf("failed to decode deposit data: %v", err)}
	}
	if err := s.store.Deposit(id, depositRequest.Amount); err != nil {
		return &shared.InternalError{Message: fmt.Sprintf("failed to deposit: %v", err)}
	}
	return WriteJSON(w, http.StatusOK, depositRequest)
}

func (s *APIServer) handeGetAllPonude(w http.ResponseWriter, _ *http.Request) error {
	ponude, err := s.get.GetAllPonude()
	if err != nil {
		return &shared.InternalError{Message: fmt.Sprintf("failed to get ponude: %v", err)}
	}
	return WriteJSON(w, http.StatusOK, ponude)
}
func getID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return 0, fmt.Errorf("invalid player id: %s", vars["id"])
	}
	return id, nil
}
