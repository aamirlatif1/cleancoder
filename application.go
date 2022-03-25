package cleancoder

import (
	"github.com/aamirlatif1/cleancoder/entity"
	"github.com/aamirlatif1/cleancoder/persistance"
)

const (
	Viewing     int8 = 1
	Downloading      = 2
)

type application struct {
	userGateway     persistance.UserGateway
	codecastGateway persistance.CodecastGateway
	licenseGateway  persistance.LicenseGateway
}

type PresentableCodecast struct {
	Title           string
	Description     string
	Picture         string
	PublicationDate string
	IsViewable      bool
	IsDownloadable  bool
}

type Gatekeeper struct {
	loggedInUser entity.User
}

func (g *Gatekeeper) SetLoggedInUser(user *entity.User) {
	g.loggedInUser = *user
}
