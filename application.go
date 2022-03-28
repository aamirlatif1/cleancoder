package cleancoder

import (
	"github.com/aamirlatif1/cleancoder/entity"
	"github.com/aamirlatif1/cleancoder/persistance"
)

type application struct {
	userGateway     persistance.UserGateway
	codecastGateway persistance.CodecastGateway
	licenseGateway  persistance.LicenseGateway
}

type Gatekeeper struct {
	loggedInUser entity.User
}

func (g *Gatekeeper) SetLoggedInUser(user *entity.User) {
	g.loggedInUser = *user
}
