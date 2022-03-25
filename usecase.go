package cleancoder

const dateLayout = "01/02/2006"

type PresentCodecastUsecase struct {
	gateway Gateway
}

func (u *PresentCodecastUsecase) PreentCodecasts(user *User) []PresentableCodecast {
	pccs := []PresentableCodecast{}
	allCodecasts := u.gateway.FindAllCodecastsSortedChronologically()

	for _, codecast := range allCodecasts {
		pccs = append(pccs, u.formatCodecast(user, codecast))
	}
	return pccs
}

func (u *PresentCodecastUsecase) formatCodecast(user *User, codecast Codecast) PresentableCodecast {
	return PresentableCodecast{
		Title:           codecast.Title,
		Description:     codecast.Title,
		Picture:         codecast.Title,
		PublicationDate: codecast.PublicationDate.Format(dateLayout),
		IsViewable:      u.IsLicenseToViewCodecast(user, &codecast),
		IsDownloadable:  u.IsLicenseToDownloadCodecast(user, &codecast),
	}
}

func (u *PresentCodecastUsecase) IsLicenseToViewCodecast(user *User, codecast *Codecast) bool {
	return u.licenseFor(user, codecast, Viewing)
}

func (u *PresentCodecastUsecase) IsLicenseToDownloadCodecast(user *User, codecast *Codecast) bool {
	return u.licenseFor(user, codecast, Downloading)
}

func (u *PresentCodecastUsecase) licenseFor(user *User, codecast *Codecast, licenseType int8) bool {
	licenses := u.gateway.FindLicensesForUserAndCodecast(user, codecast)
	for _, license := range licenses {
		if license.Type == licenseType {
			return true
		}
	}
	return false
}
