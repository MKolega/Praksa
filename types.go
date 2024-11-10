package main

type Lige struct {
	Naziv   string    `json:"naziv"`
	Razrade []Razrade `json:"razrade"`
}

type Razrade struct {
	Tipovi []Tipovi `json:"tipovi"`
	Ponude []int    `json:"ponude"`
}
type Tipovi struct {
	Naziv string `json:"naziv"`
}

type Ponude struct {
	Broj          string     `json:"broj"`
	ID            int        `json:"id"`
	Naziv         string     `json:"naziv"`
	Vrijeme       string     `json:"vrijeme"`
	Tecajevi      []Tecajevi `json:"tecajevi"`
	TvKanal       string     `json:"tv_kanal,omitempty"`
	ImaStatistiku bool       `json:"ima_statistiku,omitempty"`
}

type Tecajevi struct {
	Tecaj float64 `json:"tecaj"`
	Naziv string  `json:"naziv"`
}
type Player struct {
	ID             int     `json:"id"`
	Username       string  `json:"username"`
	Password       string  `json:"password"`
	AccountBalance float64 `json:"account_balance"`
}

type CreatePlayerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreatePonudaRequest struct {
	Broj          string     `json:"broj"`
	ID            int        `json:"id"`
	Naziv         string     `json:"naziv"`
	Vrijeme       string     `json:"vrijeme"`
	Tecajevi      []Tecajevi `json:"tecajevi"`
	TvKanal       string     `json:"tv_kanal,omitempty"`
	ImaStatistiku bool       `json:"ima_statistiku,omitempty"`
}

type OdigraniPar struct {
	Ponuda    int    `json:"ponuda"`
	NazivTipa string `json:"naziv"`
}

type DepositRequest struct {
	Amount float64 `json:"amount"`
}
type CreateUplataRequest struct {
	Amount      float64       `json:"amount"`
	OdigraniPar []OdigraniPar `json:"odigrani_par"`
}

type LigePonude struct {
	LigaNaziv string           `json:"LigaNaziv"`
	Tipovi    []string         `json:"tipovi"`
	Ponude    []PonudeTecajevi `json:"ponude"`
}

type PonudeTecajevi struct {
	NazivPonude string    `json:"NazivPonude"`
	Tecajevi    []float64 `json:"Tecajevi"`
}

func NewPonuda(broj string, ID int, naziv string, vrijeme string, tvKanal string, imaStatistiku bool) *Ponude {
	return &Ponude{
		Broj:          broj,
		ID:            ID,
		Naziv:         naziv,
		Vrijeme:       vrijeme,
		TvKanal:       tvKanal,
		ImaStatistiku: imaStatistiku,
	}
}
func NewPlayer(username, password string) *Player {
	return &Player{
		Username: username,
		Password: password,
	}
}
