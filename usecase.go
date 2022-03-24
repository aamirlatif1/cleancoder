package cleancoder

type PresentCodecastUsecase struct {
	gateway Gateway
}

func (u PresentCodecastUsecase) PreentCodecasts() []PresentableCodecast {
	return []PresentableCodecast{}
}

func (u PresentCodecastUsecase) IsLicenseToviewCodecast(user *User, codecast *CodeCast) bool {
	licenses := u.gateway.FindLicensesForUserAndCodecast(user, codecast)
	return len(licenses) > 0
}
