package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"strconv"
)

type Storage interface {
	CreatePonuda(*Ponude) error
	CreateTecaj(ponudaID int, tecaj float64, naziv string) error
	GetPonuda(id int) (*Ponude, error)
	GetAllPonude() ([]*Ponude, error)
	GetLige() ([]*Lige, error)
	CreatePlayer(*Player) error
	GetPlayers() ([]*Player, error)
	GetPlayerByID(id int) (*Player, error)
	GetLogin(username string) (*Player, error)
	Deposit(id int, amount float64) error
	CreateUplata(playerID int, amount float64, odigraniPar []OdigraniPar) error
}

type PostGresStore struct {
	db *sql.DB
}

func NewPostGresStore() (*PostGresStore, error) {
	connStr := "user=postgres dbname=postgres password=6567 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostGresStore{
		db: db,
	}, nil
}
func (s *PostGresStore) Init() error {
	return s.createPlayerTable()
}
func (s *PostGresStore) createPlayerTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS Player (
		    id SERIAL PRIMARY KEY,	
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			account_balance DECIMAL(10, 2) NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS player_ponude (
			id SERIAL PRIMARY KEY,
			player_id INT NOT NULL,
			ponuda_id INT NOT NULL,
			tecaj       NUMERIC(5, 2) NOT NULL,
			iznos_uloga NUMERIC(5, 2) NOT NULL,
			tip VARCHAR(10) NOT NULL,
			FOREIGN KEY (player_id) REFERENCES player(id) ON DELETE CASCADE,
			FOREIGN KEY (ponuda_id) REFERENCES ponude(id) ON DELETE CASCADE
		)
	`)
	return err
}

func (s *PostGresStore) CreatePlayer(player *Player) error {
	query := "INSERT INTO Player (username, password, account_balance) VALUES ($1, $2, $3)"
	resp, err := s.db.Query(query,
		player.Username, player.Password, player.AccountBalance)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *PostGresStore) GetLige() ([]*Lige, error) {
	rows, err := s.db.Query(`
		SELECT l.naziv AS liga_naziv,
		       r.ponude AS ponude,
		       t.naziv AS tip_naziv
		FROM lige l
		LEFT JOIN razrade r ON l.id = r.lige_id
		LEFT JOIN tipovi t ON r.id = t.razrade_id
		ORDER BY l.naziv, r.id, t.naziv
	`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("Failed to close rows")
		}
	}(rows)

	ligaMap := make(map[string]*Lige)

	for rows.Next() {
		var ligaNaziv string
		var ponude pq.Int64Array
		var tipNaziv sql.NullString

		if err := rows.Scan(&ligaNaziv, &ponude, &tipNaziv); err != nil {
			return nil, err
		}

		// Initialize liga
		if _, exists := ligaMap[ligaNaziv]; !exists {
			ligaMap[ligaNaziv] = &Lige{
				Naziv:   ligaNaziv,
				Razrade: []Razrade{{Tipovi: []Tipovi{}, Ponude: []int{}}},
			}
		}

		liga := ligaMap[ligaNaziv]

		// Initialize razrade
		if len(liga.Razrade) == 0 {
			liga.Razrade = append(liga.Razrade, Razrade{Tipovi: []Tipovi{}, Ponude: []int{}})
		}

		// Add tipovi
		if tipNaziv.Valid {
			tip := Tipovi{Naziv: tipNaziv.String}
			liga.Razrade[0].Tipovi = append(liga.Razrade[0].Tipovi, tip) // Add tip to the first Razrada
		}

		// Add ponude
		if len(liga.Razrade[0].Ponude) == 0 {
			for _, p := range ponude {
				liga.Razrade[0].Ponude = append(liga.Razrade[0].Ponude, int(p)) // Add ponude to the first Razrada
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var ligas []*Lige
	for _, liga := range ligaMap {
		ligas = append(ligas, liga) // Add Liga to the final slice
	}

	return ligas, nil
}

func (s *PostGresStore) CreatePonuda(ponude *Ponude) error {
	query := "INSERT INTO ponude (broj,id ,naziv,tv_kanal,vrijeme,ima_statistiku) VALUES ($1, $2, $3, $4, $5, $6)"
	resp, err := s.db.Query(query,
		ponude.Broj,
		ponude.ID,
		ponude.Naziv,
		ponude.TvKanal,
		ponude.Vrijeme,
		ponude.ImaStatistiku,
	)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *PostGresStore) CreateTecaj(ponudaID int, tecaj float64, naziv string) error {
	query := " INSERT INTO tecajevi (ponuda_id, tecaj, naziv) VALUES ($1, $2, $3)"
	resp, err := s.db.Query(query,
		ponudaID,
		tecaj,
		naziv)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *PostGresStore) GetPonuda(id int) (*Ponude, error) {
	rows, err := s.db.Query(`SELECT p.id, p.broj, p.naziv, p.vrijeme, p.tv_kanal, p.ima_statistiku, t.tecaj, t.naziv FROM ponude p LEFT JOIN tecajevi t ON p.id = t.ponuda_id WHERE p.id = $1`, id)
	if err != nil {
		return nil, err
	}
	ponuda := new(Ponude)
	ponuda.Tecajevi = []Tecajevi{}
	for rows.Next() {
		var tecaj Tecajevi
		err := rows.Scan(
			&ponuda.ID,
			&ponuda.Broj,
			&ponuda.Naziv,
			&ponuda.Vrijeme,
			&ponuda.TvKanal,
			&ponuda.ImaStatistiku,
			&tecaj.Tecaj,
			&tecaj.Naziv,
		)
		if err != nil {
			return nil, err
		}
		ponuda.Tecajevi = append(ponuda.Tecajevi, tecaj)

	}
	if ponuda.ID == 0 {
		return nil, fmt.Errorf("ponuda with id %d not found", id)
	}
	return ponuda, nil
}

func (s *PostGresStore) GetAllPonude() ([]*Ponude, error) {
	rows, err := s.db.Query(`
		SELECT p.id, p.broj, p.naziv, p.vrijeme, p.tv_kanal, p.ima_statistiku, t.tecaj, t.naziv 
		FROM ponude p 
		LEFT JOIN tecajevi t ON p.id = t.ponuda_id
		ORDER BY p.vrijeme DESC
	`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("Failed to close rows")
		}
	}(rows)

	ponudeMap := make(map[string]*Ponude)

	for rows.Next() {
		var id int
		var tecaj Tecajevi
		ponuda := Ponude{}

		err := rows.Scan(
			&id,
			&ponuda.Broj,
			&ponuda.Naziv,
			&ponuda.Vrijeme,
			&ponuda.TvKanal,
			&ponuda.ImaStatistiku,
			&tecaj.Tecaj,
			&tecaj.Naziv,
		)
		if err != nil {
			return nil, err
		}

		if existingPonuda, exists := ponudeMap[strconv.Itoa(id)]; exists {

			existingPonuda.Tecajevi = append(existingPonuda.Tecajevi, tecaj)
		} else {

			ponuda.ID = id
			ponuda.Tecajevi = []Tecajevi{tecaj}
			ponudeMap[strconv.Itoa(id)] = &ponuda
		}
	}

	ponude := make([]*Ponude, 0, len(ponudeMap))
	for _, ponuda := range ponudeMap {
		ponude = append(ponude, ponuda)
	}

	return ponude, nil
}

func (s *PostGresStore) GetPlayers() ([]*Player, error) {

	rows, err := s.db.Query(`SELECT * FROM Player`)
	if err != nil {
		return nil, err

	}
	var players []*Player
	for rows.Next() {
		player, err := scanIntoPlayer(rows)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil

}

func (s *PostGresStore) GetLogin(username string) (*Player, error) {
	rows, err := s.db.Query(`SELECT * FROM Player WHERE username = $1`, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoPlayer(rows)
	}
	return nil, fmt.Errorf("player with username %s not found", username)

}

func (s *PostGresStore) GetPlayerByID(id int) (*Player, error) {
	rows, err := s.db.Query(`SELECT * FROM Player WHERE id = $1`, id)
	if err != nil {
		return nil, err

	}

	for rows.Next() {
		return scanIntoPlayer(rows)
	}
	return nil, fmt.Errorf("player with id %d not found", id)

}

func scanIntoPlayer(rows *sql.Rows) (*Player, error) {
	player := new(Player)
	err := rows.Scan(
		&player.ID,
		&player.Username,
		&player.Password,
		&player.AccountBalance)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (s *PostGresStore) Deposit(id int, amount float64) error {
	_, err := s.db.Exec(`UPDATE player SET account_balance = account_balance + $1 WHERE id = $2`, amount, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostGresStore) CreateUplata(playerID int, amount float64, parovi []OdigraniPar) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			fmt.Println("Failed to rollback transaction")
		}
	}(tx) // Rollback in case of error

	var currentBalance float64
	err = tx.QueryRow(`SELECT account_balance FROM player WHERE id = $1`, playerID).Scan(&currentBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("player with ID %d does not exist", playerID)
		}
		return err
	}
	if currentBalance < amount {
		return fmt.Errorf("player with ID %d does not have enough funds", playerID)
	}

	for _, par := range parovi {
		var ponudaID int
		err = tx.QueryRow(`SELECT id FROM ponude WHERE id = $1`, par.Ponuda).Scan(&ponudaID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("ponuda with ID %d does not exist", par.Ponuda)
			}
		}

		var tecaj float64
		err = tx.QueryRow(`SELECT tecaj FROM tecajevi WHERE ponuda_id = $1 AND naziv = $2`, par.Ponuda, par.NazivTipa).Scan(&tecaj)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("tecaj for ponuda with ID %d and tip %s does not exist", par.Ponuda, par.NazivTipa)
			}
		}

		if amount*tecaj > 1000 {
			return fmt.Errorf("winning amount is over 1000€")
		}

		_, err = tx.Exec(`INSERT INTO player_ponude (player_id, ponuda_id, tip,tecaj,iznos_uloga) VALUES ($1, $2, $3, $4, $5)`, playerID, par.Ponuda, par.NazivTipa, tecaj, amount)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
