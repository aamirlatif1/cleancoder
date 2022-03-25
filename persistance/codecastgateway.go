package persistance

import (
	"github.com/aamirlatif1/cleancoder/entity"
)

type CodecastGateway interface {
	FindAllCodecastsSortedChronologically() []entity.Codecast
	Delete(codecast *entity.Codecast) error
	Save(codecast *entity.Codecast) *entity.Codecast
	FindCodecastByTitle(title string) (*entity.Codecast, error)
}
