package shared

type Storage interface {
	CreatePonuda(*Ponude) error
	CreateTecaj(ponudaID int, tecaj float64, naziv string) error
	GetPonuda(id int) (*Ponude, error)
	GetAllPonude() ([]*Ponude, error)
	CreateLiga(naziv string) (int, error)
	CreateRazrada(ligaID int, ponude []int) (int, error)
	CreateTipovi(razradaID int, naziv string) error
	GetLige() ([]*Lige, error)
	CreatePlayer(*Player) error
	GetPlayers() ([]*Player, error)
	GetPlayerByID(id int) (*Player, error)
	GetLogin(username string) (*Player, error)
	ResetPassword(username string, newPassword string) error
	DeleteUser(id int) error
	Deposit(id int, amount float64) error
	CreateUplata(playerID int, amount float64, odigraniPar []OdigraniPar) error
	GetAccountBalance(id int) (float64, error)
	GetPonudaByID(id int) (*Ponude, error)
	GetTecaj(parovi []OdigraniPar) ([]*OdigraniPar, error)
}

type UserError struct {
	Message string
}

func (e *UserError) Error() string {
	return e.Message
}

type InternalError struct {
	Message string
}

func (e *InternalError) Error() string {
	return e.Message
}

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

type JsonData struct {
	Lige []Lige `json:"lige"`
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
