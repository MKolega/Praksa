package main

import "time"

type Lige struct {
	Naziv   string  `json:"naziv"`
	Razrade Razrade `json:"razrade"`
	Tipovi  Tipovi  `json:"tipovi"`
}

type Razrade struct {
	Tipovi Tipovi `json:"tipovi"`
	Ponude int    `json:"ponude"`
}
type Tipovi struct {
	Naziv string `json:"naziv"`
}

type Ponude struct {
	Broj          string    `json:"broj"`
	ID            int       `json:"id"`
	Naziv         string    `json:"naziv"`
	Vrijeme       time.Time `json:"vrijeme"`
	Tecajevi      Tecajevi  `json:"tecajevi"`
	TvKanal       string    `json:"tv_kanal,omitempty"`
	ImaStatistiku bool      `json:"ima_statistiku,omitempty"`
}

type Tecajevi struct {
	Tecaj float64 `json:"tecaj"`
	Naziv string  `json:"naziv"`
}
type Player struct {
	Username       string  `json:"username"`
	Password       string  `json:"password"`
	accountBalance float64 `json:"account_balance"`
}

type CreatePlayerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewPlayer(username, password string) *Player {
	return &Player{
		Username: username,
		Password: password,
	}
}
