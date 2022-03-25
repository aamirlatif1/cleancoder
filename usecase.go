package cleancoder

type PresentCodecastUsecase struct {
	gateway Gateway
}

func (u *PresentCodecastUsecase) PreentCodecasts(user *User) []PresentableCodecast {
	presentableCodecasts := []PresentableCodecast{}
	allCodecasts := u.gateway.FindAllCodecasts()
	for _, codecast := range allCodecasts {
		presentableCodecasts = append(presentableCodecasts, PresentableCodecast{
			Title:           codecast.Title,
			Description:     codecast.Title,
			Picture:         codecast.Title,
			PublicationDate: codecast.PublicationDate,
			IsViewable:      u.IsLicenseToviewCodecast(user, &codecast),
			IsDownloadable:  false,
		})
	}
	return presentableCodecasts
}

func (u *PresentCodecastUsecase) IsLicenseToviewCodecast(user *User, codecast *CodeCast) bool {
	licenses := u.gateway.FindLicensesForUserAndCodecast(user, codecast)
	return len(licenses) > 0
}
