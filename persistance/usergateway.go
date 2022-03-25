package persistance

import (
	"github.com/aamirlatif1/cleancoder/entity"
)

type UserGateway interface {
	SaveUser(user *entity.User) *entity.User
	FindUser(username string) (*entity.User, error)
}
