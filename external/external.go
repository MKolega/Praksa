package external

import (
	"database/sql"
	"fmt"
	"github.com/MKolega/Praksa/internal/shared"
	"github.com/lib/pq"
	"strconv"
)

type PostGresStore struct {
	db *sql.DB
}

func NewPostGresGet() (*PostGresStore, error) {
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

func (s *PostGresStore) GetLige() ([]*shared.Lige, error) {
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

	ligaMap := make(map[string]*shared.Lige)

	for rows.Next() {
		var ligaNaziv string
		var ponude pq.Int64Array
		var tipNaziv sql.NullString

		if err := rows.Scan(&ligaNaziv, &ponude, &tipNaziv); err != nil {
			return nil, err
		}

		// Initialize liga
		if _, exists := ligaMap[ligaNaziv]; !exists {
			ligaMap[ligaNaziv] = &shared.Lige{
				Naziv:   ligaNaziv,
				Razrade: []shared.Razrade{{Tipovi: []shared.Tipovi{}, Ponude: []int{}}},
			}
		}

		liga := ligaMap[ligaNaziv]

		// Initialize razrade
		if len(liga.Razrade) == 0 {
			liga.Razrade = append(liga.Razrade, shared.Razrade{Tipovi: []shared.Tipovi{}, Ponude: []int{}})
		}

		// Add tipovi
		if tipNaziv.Valid {
			tip := shared.Tipovi{Naziv: tipNaziv.String}
			liga.Razrade[0].Tipovi = append(liga.Razrade[0].Tipovi, tip)
		}

		// Add ponude
		if len(liga.Razrade[0].Ponude) == 0 {
			for _, p := range ponude {
				liga.Razrade[0].Ponude = append(liga.Razrade[0].Ponude, int(p))
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var ligas []*shared.Lige
	for _, liga := range ligaMap {
		ligas = append(ligas, liga)
	}

	return ligas, nil
}

func (s *PostGresStore) GetPonuda(id int) (*shared.Ponude, error) {
	rows, err := s.db.Query(`SELECT p.id, p.broj, p.naziv, p.vrijeme, p.tv_kanal, p.ima_statistiku, t.tecaj, t.naziv FROM ponude p LEFT JOIN tecajevi t ON p.id = t.ponuda_id WHERE p.id = $1`, id)
	if err != nil {
		return nil, err
	}
	ponuda := new(shared.Ponude)
	ponuda.Tecajevi = []shared.Tecajevi{}
	for rows.Next() {
		var tecaj shared.Tecajevi
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

func (s *PostGresStore) GetAllPonude() ([]*shared.Ponude, error) {
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

	ponudeMap := make(map[string]*shared.Ponude)

	for rows.Next() {
		var id int
		var tecaj shared.Tecajevi
		ponuda := shared.Ponude{}

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
			ponuda.Tecajevi = []shared.Tecajevi{tecaj}
			ponudeMap[strconv.Itoa(id)] = &ponuda
		}
	}

	ponude := make([]*shared.Ponude, 0, len(ponudeMap))
	for _, ponuda := range ponudeMap {
		ponude = append(ponude, ponuda)
	}

	return ponude, nil
}

func (s *PostGresStore) GetPlayers() ([]*shared.Player, error) {

	rows, err := s.db.Query(`SELECT * FROM Player`)
	if err != nil {
		return nil, err

	}
	var players []*shared.Player
	for rows.Next() {
		player, err := scanIntoPlayer(rows)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil

}

func (s *PostGresStore) GetLogin(username string) (*shared.Player, error) {
	rows, err := s.db.Query(`SELECT * FROM Player WHERE username = $1`, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoPlayer(rows)
	}
	return nil, fmt.Errorf("player with username %s not found", username)

}

func (s *PostGresStore) GetPlayerByID(id int) (*shared.Player, error) {
	rows, err := s.db.Query(`SELECT * FROM Player WHERE id = $1`, id)
	if err != nil {
		return nil, err

	}

	for rows.Next() {
		return scanIntoPlayer(rows)
	}
	return nil, fmt.Errorf("player with id %d not found", id)

}

func scanIntoPlayer(rows *sql.Rows) (*shared.Player, error) {
	player := new(shared.Player)
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
