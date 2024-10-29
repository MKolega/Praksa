package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreatePonuda(*Ponude) error
	GetPonuda(id int) (*Ponude, error)
	GetLige() ([]*Lige, error)
	CreatePlayer(*Player) error
	GetPlayers() ([]*Player, error)
	GetPlayerByID(id int) (*Player, error)
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
		)
	`)
	return err
}

func (s *PostGresStore) CreatePlayer(player *Player) error {
	query := "INSERT INTO Player (username, password, account_balance) VALUES ($1, $2, $3)"
	resp, err := s.db.Query(query,
		player.Username, player.Password, player.accountBalance)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *PostGresStore) GetLige() ([]*Lige, error) {
	rows, err := s.db.Query(`SELECT naziv,ponuda_id FROM lige,lige_ponude WHERE lige.id=lige_ponude.lige_id`)
	if err != nil {
		return nil, err
	}

	lige := []*Lige{}
	for rows.Next() {
		liga := new(Lige)
		err := rows.Scan(
			&liga.Naziv,
			&liga.Razrade.Ponude,
		)
		if err != nil {
			return nil, err
		}
		lige = append(lige, liga)

		/*
			if err := json.Unmarshal([]byte(razradeJSON), &liga.Tipovi); err != nil {
				return nil, err
			}*/
	}
	return lige, nil
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

func (s *PostGresStore) GetPonuda(id int) (*Ponude, error) {
	rows, err := s.db.Query(`SELECT id, broj, naziv, vrijeme, tv_kanal, ima_statistiku FROM ponude WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	ponuda := new(Ponude)
	for rows.Next() {
		err := rows.Scan(
			&ponuda.ID,
			&ponuda.Broj,
			&ponuda.Naziv,
			&ponuda.Vrijeme,
			&ponuda.TvKanal,
			&ponuda.ImaStatistiku,
		)
		if err != nil {
			return nil, err
		}

	}
	if ponuda.ID == 0 {
		return nil, fmt.Errorf("ponuda with id %d not found", id)
	}
	return ponuda, nil
}

func (s *PostGresStore) GetPlayers() ([]*Player, error) {

	rows, err := s.db.Query(`SELECT * FROM Player`)
	if err != nil {
		return nil, err

	}
	players := []*Player{}
	for rows.Next() {
		player, err := scanIntoPlayer(rows)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil

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
		&player.accountBalance)
	if err != nil {
		return nil, err
	}
	return player, nil
}

/*
func (s *PostGresStore) CreatePonuda(ponude Ponude) error {
	return nil
}
*/
/*
	func (s *PostGresStore) GetPonuda(id int) ([]Ponude, error) {
		ponuda := Ponude{}

		query := "SELECT broj, id, naziv, vrijeme, tv_kanal, ima_statistiku FROM ponude WHERE id = $1"
		row := s.db.QueryRow(query, id)
		err := row.Scan(&ponuda.Broj, &ponuda.ID, &ponuda.Naziv, &ponuda.Vrijeme, &ponuda.TvKanal, &ponuda.ImaStatistiku)
		if err != nil {
			return Ponuda{}, err
		}
		return ponuda, nil
	}
*/
/*func (s *PostGresStore) GetLige() (Lige, error) {
	rows, err := s.db.Query(`SELECT naziv,ponude FROM lige`)
	if err != nil {
		return Lige{}, err
	}
	lige := Lige{}
	for rows.Next() {
		liga := new(Lige)
		err := rows.Scan(&liga.Naziv, &liga.Razrade)
		if err != nil {
			return Lige{}, err
		}
		fmt.Printf("Liga: %s\n", liga.Naziv)
		for _, razrada := range liga.Razrade {
			for _, ponuda := range razrada.Ponude {
				fmt.Printf("Ponuda: %d\n", ponuda)
			}
		}

	}
	return lige, nil
}
*/
