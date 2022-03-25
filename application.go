package cleancoder

type application struct {
	gateway Gateway
}

type Entity interface {
	SetID(id string)
	GetID() string
	IsSame(entity Entity) bool
}

type CodeCast struct {
	ID              string
	Title           string
	PublicationDate string
}

func (c *CodeCast) SetID(id string) {
	c.ID = id
}

func (c *CodeCast) GetID() string {
	return c.ID
}

func (c *CodeCast) IsSame(entity Entity) bool {
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
	FindAllCodecasts() []CodeCast
	Delete(codecast *CodeCast) error
	Save(codecast *CodeCast) *CodeCast
	SaveUser(user *User) *User
	FindUser(username string) (*User, error)
	FindCodecastByTitle(title string) (*CodeCast, error)
	SaveLicense(liccense *License)
	FindLicensesForUserAndCodecast(user *User, codecast *CodeCast) []License
}

type Gatekeeper struct {
	loggedInUser User
}

func (g *Gatekeeper) SetLoggedInUser(user *User) {
	g.loggedInUser = *user
}

type License struct {
	User     User
	codecast CodeCast
}
