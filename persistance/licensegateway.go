package persistance

import (
	"github.com/aamirlatif1/cleancoder/entity"
)

type LicenseGateway interface {
	SaveLicense(liccense *entity.License)
	FindLicensesForUserAndCodecast(user *entity.User, codecast *entity.Codecast) []entity.License
}
