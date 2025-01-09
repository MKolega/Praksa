package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MKolega/Praksa/internal/shared"
	"github.com/lib/pq"
	"log"
)

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
		);
		
		CREATE TABLE IF NOT EXISTS lige (
		    			id SERIAL PRIMARY KEY,
		    			naziv VARCHAR(255) NOT NULL
		                                		    		);

		CREATE TABLE IF NOT EXISTS razrade (
		    			id SERIAL PRIMARY KEY,
		    			lige_id INT NOT NULL,
		    			ponude INT[] NOT NULL,
		    			FOREIGN KEY (lige_id) REFERENCES lige(id) ON DELETE CASCADE
		                                   		                                		    		);

		CREATE TABLE IF NOT EXISTS tipovi (
		    			id SERIAL PRIMARY KEY,
		    			razrade_id INT NOT NULL,
		    			naziv VARCHAR(255) NOT NULL,
		    			FOREIGN KEY (razrade_id) REFERENCES razrade(id) ON DELETE CASCADE
		                                   		                                		    		);

		CREATE TABLE IF NOT EXISTS ponude (
		    			id SERIAL PRIMARY KEY,
		    			broj VARCHAR(255) NOT NULL,
		    			naziv VARCHAR(255) NOT NULL,
		    			tv_kanal VARCHAR(255) default NULL,
		    			vrijeme TIMESTAMP NOT NULL,
		    			ima_statistiku BOOLEAN default FALSE
		                                  		                                		    		);

		CREATE TABLE IF NOT EXISTS tecajevi (
		    			id SERIAL PRIMARY KEY,
		    			ponuda_id INT NOT NULL,
		    			tecaj NUMERIC(5, 2) NOT NULL,
		    			naziv VARCHAR(255) NOT NULL,
		    			FOREIGN KEY (ponuda_id) REFERENCES ponude(id) ON DELETE CASCADE
		                                    		                                   		                                		    		);


	`)
	return err
}

func (s *PostGresStore) CreatePlayer(player *shared.Player) error {
	query := "INSERT INTO Player (username, password, account_balance) VALUES ($1, $2, $3)"
	resp, err := s.db.Query(query,
		player.Username, player.Password, player.AccountBalance)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *PostGresStore) CreateLiga(naziv string) (int, error) {
	var existingID int
	checkQuery := `SELECT id FROM lige WHERE naziv = $1`
	err := s.db.QueryRow(checkQuery, naziv).Scan(&existingID)
	if err == nil {
		log.Printf("Duplicate liga found: %s (ID: %d), skipping insertion.", naziv, existingID)
		return 0, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("failed to check for duplicate liga: %v", err)
	}

	query := "INSERT INTO ligetest (naziv) VALUES ($1) RETURNING id"
	var ligaID int
	err = s.db.QueryRow(query, naziv).Scan(&ligaID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert liga: %v", err)
	}
	return ligaID, nil
}

func (s *PostGresStore) CreateRazrada(ligaID int, ponude []int) (int, error) {
	var razradaID int
	err := s.db.QueryRow(`INSERT INTO razrade (lige_id, ponude) VALUES ($1, $2) RETURNING id`,
		ligaID, pq.Array(ponude)).Scan(&razradaID)
	if err != nil {
		return 0, fmt.Errorf("failed to create razrada: %v", err)
	}
	return razradaID, nil
}

func (s *PostGresStore) CreateTipovi(razradaID int, naziv string) error {
	query := `INSERT INTO tipovi (razrade_id, naziv) VALUES ($1, $2)`
	_, err := s.db.Exec(query, razradaID, naziv)
	if err != nil {
		return fmt.Errorf("failed to insert tip: %v", err)
	}

	return nil
}

func (s *PostGresStore) CreatePonuda(ponude *shared.Ponude) error {
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
		return fmt.Errorf("failed to insert ponude: %v", err)
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
		return fmt.Errorf("failed to insert tecaj: %v", err)
	}
	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *PostGresStore) ResetPassword(username string, newPassword string) error {
	_, err := s.db.Exec(`UPDATE Player SET password = $1 WHERE username = $2`, newPassword, username)
	if err != nil {
		return err
	}
	return nil

}

func (s *PostGresStore) DeleteUser(id int) error {
	_, err := s.db.Exec(`
    DELETE FROM player_ponude WHERE player_id = $1;
    DELETE FROM Player WHERE id = $1;
`, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostGresStore) Deposit(id int, amount float64) error {
	_, err := s.db.Exec(`UPDATE player SET account_balance = account_balance + $1 WHERE id = $2`, amount, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostGresStore) CreateUplata(playerID int, amount float64, parovi []shared.OdigraniPar) error {
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
