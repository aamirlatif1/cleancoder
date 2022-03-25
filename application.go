package cleancoder

import "time"

const (
	Viewing     int8 = 1
	Downloading      = 2
)

type application struct {
	gateway Gateway
}

type Entity interface {
	SetID(id string)
	GetID() string
	IsSame(entity Entity) bool
}

type Codecast struct {
	ID              string
	Title           string
	PublicationDate time.Time
}

func (c *Codecast) SetID(id string) {
	c.ID = id
}

func (c *Codecast) GetID() string {
	return c.ID
}

func (c *Codecast) IsSame(entity Entity) bool {
	return c.ID != "" && c.ID == entity.GetID()
}

type PresentableCodecast struct {
	Title           string
	Description     string
	Picture         string
	PublicationDate string
	IsViewable      bool
	IsDownloadable  bool
}

type User struct {
	Username string
	ID       string
}

func (u *User) SetID(id string) {
	u.ID = id
}

func (u *User) GetID() string {
	return u.ID
}

func (u *User) IsSame(entity Entity) bool {
	return u.ID != "" && u.ID == entity.GetID()
}

type Gateway interface {
	FindAllCodecastsSortedChronologically() []Codecast
	Delete(codecast *Codecast) error
	Save(codecast *Codecast) *Codecast
	SaveUser(user *User) *User
	FindUser(username string) (*User, error)
	FindCodecastByTitle(title string) (*Codecast, error)
	SaveLicense(liccense *License)
	FindLicensesForUserAndCodecast(user *User, codecast *Codecast) []License
}

type Gatekeeper struct {
	loggedInUser User
}

func (g *Gatekeeper) SetLoggedInUser(user *User) {
	g.loggedInUser = *user
}

type License struct {
	User     User
	Codecast Codecast
	Type     int8
}
