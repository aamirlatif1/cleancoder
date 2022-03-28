package usecase

import (
	"github.com/aamirlatif1/cleancoder/data"
	"github.com/aamirlatif1/cleancoder/entity"
	"github.com/aamirlatif1/cleancoder/persistance"
)

const dateLayout = "01/02/2006"

type PresentCodecast struct {
	UserGateway     persistance.UserGateway
	CodecastGateway persistance.CodecastGateway
	LicenseGateway  persistance.LicenseGateway
}

func (u *PresentCodecast) PreentCodecasts(user *entity.User) []data.PresentableCodecast {
	pccs := []data.PresentableCodecast{}
	allCodecasts := u.CodecastGateway.FindAllCodecastsSortedChronologically()

	for _, codecast := range allCodecasts {
		pccs = append(pccs, u.formatCodecast(user, codecast))
	}
	return pccs
}

func (u *PresentCodecast) formatCodecast(user *entity.User, codecast entity.Codecast) data.PresentableCodecast {
	return data.PresentableCodecast{
		Title:           codecast.Title,
		Description:     codecast.Title,
		Picture:         codecast.Title,
		PublicationDate: codecast.PublicationDate.Format(dateLayout),
		IsViewable:      u.IsLicenseToViewCodecast(user, &codecast),
		IsDownloadable:  u.IsLicenseToDownloadCodecast(user, &codecast),
	}
}

func (u *PresentCodecast) IsLicenseToViewCodecast(user *entity.User, codecast *entity.Codecast) bool {
	return u.licenseFor(user, codecast, data.Viewing)
}

func (u *PresentCodecast) IsLicenseToDownloadCodecast(user *entity.User, codecast *entity.Codecast) bool {
	return u.licenseFor(user, codecast, data.Downloading)
}

func (u *PresentCodecast) licenseFor(user *entity.User, codecast *entity.Codecast, licenseType int8) bool {
	licenses := u.LicenseGateway.FindLicensesForUserAndCodecast(user, codecast)
	for _, license := range licenses {
		if license.Type == licenseType {
			return true
		}
	}
	return false
}
