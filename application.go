package cleancoder

type application struct {
	gateway Gateway
}

type CodeCast struct {
	Title string
	PublishedDate string
}

type PresentableCodecast struct {

}


type User struct {
	Username string
}

type Gateway interface {
	FindAllCodecasts() []CodeCast
	Delete(codecast *CodeCast) error
	Save(codecast *CodeCast)
	SaveUser(user *User)
	FindUser(username string) (*User, error)
}

type Gatekeeper struct {
	loggedInUser User
}

func (g *Gatekeeper) SetLoggedInUser(user *User)  {
	g.loggedInUser = *user
}


