package cleancoder

import (
	"github.com/aamirlatif1/cleancoder/entity"
	"github.com/aamirlatif1/cleancoder/persistance"
)

const dateLayout = "01/02/2006"

type PresentCodecastUsecase struct {
	userGateway     persistance.UserGateway
	codecastGateway persistance.CodecastGateway
	licenseGateway  persistance.LicenseGateway
}

func (u *PresentCodecastUsecase) PreentCodecasts(user *entity.User) []PresentableCodecast {
	pccs := []PresentableCodecast{}
	allCodecasts := u.codecastGateway.FindAllCodecastsSortedChronologically()

	for _, codecast := range allCodecasts {
		pccs = append(pccs, u.formatCodecast(user, codecast))
	}
	return pccs
}

func (u *PresentCodecastUsecase) formatCodecast(user *entity.User, codecast entity.Codecast) PresentableCodecast {
	return PresentableCodecast{
		Title:           codecast.Title,
		Description:     codecast.Title,
		Picture:         codecast.Title,
		PublicationDate: codecast.PublicationDate.Format(dateLayout),
		IsViewable:      u.IsLicenseToViewCodecast(user, &codecast),
		IsDownloadable:  u.IsLicenseToDownloadCodecast(user, &codecast),
	}
}

func (u *PresentCodecastUsecase) IsLicenseToViewCodecast(user *entity.User, codecast *entity.Codecast) bool {
	return u.licenseFor(user, codecast, Viewing)
}

func (u *PresentCodecastUsecase) IsLicenseToDownloadCodecast(user *entity.User, codecast *entity.Codecast) bool {
	return u.licenseFor(user, codecast, Downloading)
}

func (u *PresentCodecastUsecase) licenseFor(user *entity.User, codecast *entity.Codecast, licenseType int8) bool {
	licenses := u.licenseGateway.FindLicensesForUserAndCodecast(user, codecast)
	for _, license := range licenses {
		if license.Type == licenseType {
			return true
		}
	}
	return false
}
