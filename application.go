package cleancoder

type application struct {
	gateway Gateway
}

type CodeCast struct {
	Title         string
	PublishedDate string
}

func (c *CodeCast) IsSame(codecast *CodeCast) bool {
	return c.Title == codecast.Title
}

type PresentableCodecast struct {
}

type User struct {
	Username string
	ID       string
}

func (u *User) IsSame(user *User) bool {
	return u.ID == user.ID
}

type Gateway interface {
	FindAllCodecasts() []CodeCast
	Delete(codecast *CodeCast) error
	Save(codecast *CodeCast)
	SaveUser(user *User)
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
